import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { Item } from '../types';
import { BaseScraper } from './base';
import { NewsCache } from '../utils/cache';

export class PhoronixScraper extends BaseScraper {
  private baseUrl = 'https://www.phoronix.com';

  constructor(sharedCache?: NewsCache) {
    super(sharedCache);
  }

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

        const item = await this.processNewsItem(href, linkData, 'Phoronix');
        if (item) {
          items.push(item);
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