import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { Item } from '../types';
import { BaseScraper } from './base';

export class LinuxComScraper extends BaseScraper {
  private baseUrl = 'https://www.linux.com';

  async scrape(): Promise<Item[]> {
    const url = `${this.baseUrl}/news/`;
    logger.info('Iniciando scraping do Linux.com', { url });

    try {
      const html = await this.fetchPage(url);
      const $ = cheerio.load(html);
      const items: Item[] = [];

      const newsLinks = $('.post-title a, .entry-title a, h2 a, h3 a').map((i, el) => {
        const href = $(el).attr('href');
        const title = $(el).text().trim();
        
        if (href && title && title.length > 10) {
          return {
            href: href.startsWith('http') ? href : `${this.baseUrl}${href}`,
            title: title
          };
        }
        return null;
      }).get().filter(Boolean);

      const hrefs = newsLinks.map(link => link.href);
      const newLinks = this.getUniqueNewLinks(hrefs);
      
      this.logLinkStats('Linux.com', newsLinks.length, hrefs.length, newLinks.length);

      if (newLinks.length === 0) {
        return [];
      }

      for (const href of newLinks.slice(0, this.maxItems)) {
        const linkData = newsLinks.find(link => link.href === href);
        if (!linkData) continue;

        const item: Item = {
          type: 'news',
          title: linkData.title,
          author: 'Linux.com',
          date: 'Desconhecida',
          body: '',
          link: href
        };

        try {
          await this.delay(1000);
          const newsHtml = await this.fetchPage(href);
          const $news = cheerio.load(newsHtml);

          const title = $news('h1, .entry-title, .post-title').first().text().trim();
          const author = $news('.author, .byline, [rel="author"]').first().text().trim();
          const date = $news('.published, .date, time[datetime]').first().text().trim() || 
                      $news('.published, .date, time[datetime]').first().attr('datetime');
          const bodyText = $news('article p, .entry-content p, .post-content p')
                          .map((i, p) => $(p).text().trim())
                          .get()
                          .filter(text => text.length > 20)
                          .join(' ');

          item.title = title || item.title;
          item.author = author || item.author;
          item.date = date || item.date;
          item.body = bodyText || 'Conteúdo não disponível';

          this.cache.markAsSeen(href);
          items.push(item);

          logger.info('Notícia processada do Linux.com', {
            title: item.title.substring(0, 50),
            author: item.author
          });

        } catch (error) {
          logger.error('Erro ao processar notícia do Linux.com', {
            url: href,
            error: (error as Error).message
          });
        }
      }

      logger.info('Scraping do Linux.com concluído', { 
        itemsColetados: items.length 
      });

      return items;

    } catch (error) {
      logger.error('Erro no scraping do Linux.com', {
        error: (error as Error).message
      });
      return [];
    }
  }
}