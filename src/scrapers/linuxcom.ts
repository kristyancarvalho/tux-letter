import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { Item } from '../types';
import { BaseScraper } from './base';
import { NewsCache } from '../utils/cache';

export class LinuxComScraper extends BaseScraper {
  private baseUrl = 'https://www.linux.com';

  constructor(sharedCache?: NewsCache) {
    super(sharedCache);
  }

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

        const item = await this.processNewsItem(href, linkData, 'Linux.com');
        if (item) {
          items.push(item);
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