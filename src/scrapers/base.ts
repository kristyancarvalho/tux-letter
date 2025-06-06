import axios from 'axios';
import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { NewsCache } from '../utils/cache';
import { Item } from '../types';

export abstract class BaseScraper {
  protected cache: NewsCache;
  protected maxItems: number = 5;

  constructor(sharedCache?: NewsCache) {
    this.cache = sharedCache || new NewsCache();
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

  protected async processNewsItem(href: string, linkData: any, defaultAuthor: string): Promise<Item | null> {
    if (!this.cache.isNew(href)) {
      return null;
    }

    this.cache.markAsSeen(href);

    const item: Item = {
      type: 'news',
      title: linkData.title,
      author: defaultAuthor,
      date: 'Desconhecida',
      body: '',
      link: href
    };

    try {
      await this.delay(1000);
      const newsHtml = await this.fetchPage(href);
      const $news = cheerio.load(newsHtml);

      const title = $news('h1, .post-title, .entry-title').first().text().trim();
      const author = $news('.author-name, .byline, [rel="author"], .author, [class*="author"]').first().text().trim();
      const date = $news('.published-date, .post-date, time, .published, .date, [datetime]').first().text().trim() ||
                  $news('time[datetime]').first().attr('datetime');
      const bodyText = $news('article p, .post-content p, .entry-content p, .content p')
                      .map((i, p) => $news(p).text().trim())
                      .get()
                      .filter(text => text.length > 20)
                      .join(' ');
                      
      item.title = title || item.title;
      item.author = author || item.author;
      item.date = date || item.date;
      item.body = bodyText || 'Conteúdo não disponível';

      logger.info(`Notícia processada de ${defaultAuthor}`, {
        title: item.title.substring(0, 50),
        author: item.author
      });

      return item;

    } catch (error) {
      logger.error(`Erro ao processar notícia de ${defaultAuthor}`, {
        url: href,
        error: (error as Error).message
      });
      return null;
    }
  }

  abstract scrape(): Promise<Item[]>;
}