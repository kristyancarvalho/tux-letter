import { LoreScraper } from './lore';
import { PhoronixScraper } from './phoronix';
import { LinuxComScraper } from './linuxcom';
import { ItsFossScraper } from './itsfoss';
import { Item, BaseScraper } from '../types';
import { logger } from '../utils/logger';
import { NewsCache } from '../utils/cache';

export class ScraperManager {
  private scrapers: Map<string, BaseScraper>;
  private cache: NewsCache;

  constructor() {
    this.cache = new NewsCache();
    this.scrapers = new Map<string, BaseScraper>();
    
    this.scrapers.set('lore', new LoreScraper());
    this.scrapers.set('phoronix', new PhoronixScraper());
    this.scrapers.set('linuxcom', new LinuxComScraper());
    this.scrapers.set('itsfoss', new ItsFossScraper());
  }

  async scrapeAll(): Promise<Item[]> {
    logger.info('🔄 Iniciando scraping de todos os sites');
    
    const allItems: Item[] = [];
    const results = new Map<string, Item[]>();

    for (const [name, scraper] of this.scrapers) {
      try {
        logger.info(`Executando scraper: ${name}`);
        const items = await scraper.scrape();
        results.set(name, items);
        allItems.push(...items);
        
        logger.info(`Scraper ${name} concluído`, { 
          itemsEncontrados: items.length 
        });
        
      } catch (error) {
        logger.error(`Erro no scraper ${name}`, {
          error: (error as Error).message
        });
        results.set(name, []);
      }
    }

    const uniqueItems = this.removeDuplicates(allItems);
    
    logger.info('📊 Resumo do scraping completo', {
      loreItems: results.get('lore')?.length || 0,
      phoronixItems: results.get('phoronix')?.length || 0,
      linuxcomItems: results.get('linuxcom')?.length || 0,
      itsfossItems: results.get('itsfoss')?.length || 0,
      totalItems: allItems.length,
      uniqueItems: uniqueItems.length
    });

    return uniqueItems;
  }

  async scrapeSpecific(scraperNames: string[]): Promise<Item[]> {
    const allItems: Item[] = [];

    for (const name of scraperNames) {
      const scraper = this.scrapers.get(name);
      if (!scraper) {
        logger.warn(`Scraper não encontrado: ${name}`);
        continue;
      }

      try {
        const items = await scraper.scrape();
        allItems.push(...items);
      } catch (error) {
        logger.error(`Erro no scraper ${name}`, {
          error: (error as Error).message
        });
      }
    }

    return this.removeDuplicates(allItems);
  }

  private removeDuplicates(items: Item[]): Item[] {
    const seenLinks = new Set<string>();
    const uniqueItems: Item[] = [];

    items.forEach(item => {
      if (!seenLinks.has(item.link)) {
        seenLinks.add(item.link);
        uniqueItems.push(item);
      }
    });

    return uniqueItems;
  }

  getCacheStats() {
    return this.cache.getStats();
  }

  clearCache() {
    this.cache.clear();
  }

  persistCache() {
    this.cache.persist();
  }

  getBotVerificationCount(): number {
    const loreScraper = this.scrapers.get('lore');
    return (loreScraper && 'getBotVerificationCount' in loreScraper) 
      ? (loreScraper as any).getBotVerificationCount() 
      : 0;
  }
}