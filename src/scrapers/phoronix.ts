import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { Item } from '../types';
import { BaseScraper } from './base';

export class PhoronixScraper extends BaseScraper {
  private baseUrl = 'https://www.phoronix.com';

  async scrape(): Promise<Item[]> {
    const url = this.baseUrl;
    logger.info('Iniciando scraping do Phoronix', { url });

    try {
      const html = await this.fetchPage(url);
      const $ = cheerio.load(html);
      const items: Item[] = [];

      const newsLinks = $('a[href*="/news/"]').map((i, el) => {
        const href = $(el).attr('href');
        const title = $(el).text().trim();
        
        if (href && title && href.includes('/news/') && !href.includes('rss') && 
            $(el).closest('.sidebar').length === 0 && title.length > 10) {
          return {
            href: href.startsWith('http') ? href : `${this.baseUrl}${href}`,
            title: title
          };
        }
        return null;
      }).get().filter(Boolean);

      const hrefs = newsLinks.map(link => link.href);
      const newLinks = this.getUniqueNewLinks(hrefs);
      
      this.logLinkStats('Phoronix', newsLinks.length, hrefs.length, newLinks.length);

      if (newLinks.length === 0) {
        return [];
      }

      for (const href of newLinks.slice(0, this.maxItems)) {
        const linkData = newsLinks.find(link => link.href === href);
        if (!linkData) continue;

        const item: Item = {
          type: 'news',
          title: linkData.title,
          author: 'Phoronix',
          date: 'Desconhecida',
          body: '',
          link: href
        };

        try {
          await this.delay(1000);
          const newsHtml = await this.fetchPage(href);
          const $news = cheerio.load(newsHtml);

          const title = $news('h1').first().text().trim();
          const author = $news('.author, .byline, [class*="author"]').first().text().trim();
          const date = $news('.date, .published, time, [datetime]').first().text().trim();
          const bodyText = $news('article p, .content p').map((i, p) => $(p).text().trim()).get().join(' ');

          item.title = title || item.title;
          item.author = author || item.author;
          item.date = date || item.date;
          item.body = bodyText || 'Conteúdo não disponível';

          this.cache.markAsSeen(href);
          items.push(item);

          logger.info('Notícia processada do Phoronix', {
            title: item.title.substring(0, 50),
            author: item.author
          });

        } catch (error) {
          logger.error('Erro ao processar notícia do Phoronix', {
            url: href,
            error: (error as Error).message
          });
        }
      }

      logger.info('Scraping do Phoronix concluído', { 
        itemsColetados: items.length 
      });

      return items;

    } catch (error) {
      logger.error('Erro no scraping do Phoronix', {
        error: (error as Error).message
      });
      return [];
    }
  }
}