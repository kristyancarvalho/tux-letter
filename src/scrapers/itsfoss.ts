import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { Item } from '../types';
import { BaseScraper } from './base';
import { NewsCache } from '../utils/cache';

export class ItsFossScraper extends BaseScraper {
  private baseUrl = 'https://news.itsfoss.com';

  constructor(sharedCache?: NewsCache) {
    super(sharedCache);
  }

  async scrape(): Promise<Item[]> {
    const url = this.baseUrl;
    logger.info('Iniciando scraping do Its FOSS News', { url });

    try {
      const html = await this.fetchPage(url);
      const $ = cheerio.load(html);
      const items: Item[] = [];

      const newsLinks = $('.post-card-title a, .entry-title a, h2 a, h3 a').map((i, el) => {
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
      
      this.logLinkStats('Its FOSS News', newsLinks.length, hrefs.length, newLinks.length);

      if (newLinks.length === 0) {
        return [];
      }

      for (const href of newLinks.slice(0, this.maxItems)) {
        const linkData = newsLinks.find(link => link.href === href);
        if (!linkData) continue;

        const item = await this.processNewsItem(href, linkData, 'Its FOSS');
        if (item) {
          items.push(item);
        }
      }

      logger.info('Scraping do Its FOSS News concluído', { 
        itemsColetados: items.length 
      });

      return items;

    } catch (error) {
      logger.error('Erro no scraping do Its FOSS News', {
        error: (error as Error).message
      });
      return [];
    }
  }
}