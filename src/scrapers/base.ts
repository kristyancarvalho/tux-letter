import axios from 'axios';
import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { NewsCache } from '../utils/cache';
import { Item } from '../types';

export abstract class BaseScraper {
  protected cache: NewsCache;
  protected maxItems: number = 5;

  constructor() {
    this.cache = new NewsCache();
  }

  protected async fetchPage(url: string): Promise<string> {
    const response = await axios.get(url, {
      headers: {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Connection': 'keep-alive'
      },
      timeout: 30000
    });
    return response.data;
  }

  protected async delay(ms: number = 1000): Promise<void> {
    await new Promise(resolve => setTimeout(resolve, ms));
  }

  protected getUniqueNewLinks(links: string[]): string[] {
    const uniqueLinks = Array.from(new Set(links));
    return this.cache.filterNewLinks(uniqueLinks).slice(0, this.maxItems);
  }

  protected logLinkStats(siteName: string, totalLinks: number, uniqueLinks: number, newLinks: number): void {
    logger.info(`Links processados em ${siteName}`, {
      total: totalLinks,
      unique: uniqueLinks,
      new: newLinks,
      cached: uniqueLinks - newLinks
    });
  }

  abstract scrape(): Promise<Item[]>;
}