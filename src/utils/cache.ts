import fs from 'fs';
import path from 'path';
import { logger } from './logger';

interface CacheData {
  seenLinks: string[];
  lastUpdated: string;
}

export class NewsCache {
  private cacheFile: string;
  private seenLinks: Set<string> = new Set();

  constructor(cacheDir: string = path.join(__dirname, '../../cache')) {
    if (!fs.existsSync(cacheDir)) {
      fs.mkdirSync(cacheDir, { recursive: true });
    }
    
    this.cacheFile = path.join(cacheDir, 'seen_links.json');
    this.loadCache();
  }

  private loadCache(): void {
    try {
      if (fs.existsSync(this.cacheFile)) {
        const cacheData: CacheData = JSON.parse(fs.readFileSync(this.cacheFile, 'utf8'));
        this.seenLinks = new Set(cacheData.seenLinks);
        logger.info('Cache carregado', { 
          totalLinks: this.seenLinks.size,
          lastUpdated: cacheData.lastUpdated 
        });
      } else {
        logger.info('Arquivo de cache não encontrado, iniciando com cache vazio');
      }
    } catch (error) {
      logger.error('Erro ao carregar cache', { 
        error: (error as Error).message,
        cacheFile: this.cacheFile 
      });
      this.seenLinks = new Set();
    }
  }

  private saveCache(): void {
    try {
      const cacheData: CacheData = {
        seenLinks: Array.from(this.seenLinks),
        lastUpdated: new Date().toISOString()
      };
      
      fs.writeFileSync(this.cacheFile, JSON.stringify(cacheData, null, 2));
      logger.info('Cache salvo', { 
        totalLinks: this.seenLinks.size,
        cacheFile: this.cacheFile 
      });
    } catch (error) {
      logger.error('Erro ao salvar cache', { 
        error: (error as Error).message,
        cacheFile: this.cacheFile 
      });
    }
  }

  isNew(link: string): boolean {
    return !this.seenLinks.has(link);
  }

  markAsSeen(link: string): void {
    if (this.isNew(link)) {
      this.seenLinks.add(link);
      logger.debug('Link marcado como visto', { link });
    }
  }

  markMultipleAsSeen(links: string[]): void {
    let newLinksCount = 0;
    links.forEach(link => {
      if (this.isNew(link)) {
        this.seenLinks.add(link);
        newLinksCount++;
      }
    });
    
    if (newLinksCount > 0) {
      logger.info('Múltiplos links marcados como vistos', { 
        newLinks: newLinksCount,
        totalLinks: links.length 
      });
    }
  }

  filterNewLinks(links: string[]): string[] {
    const newLinks = links.filter(link => this.isNew(link));
    logger.info('Links filtrados', { 
      totalLinks: links.length,
      newLinks: newLinks.length,
      skippedLinks: links.length - newLinks.length 
    });
    return newLinks;
  }

  persist(): void {
    this.saveCache();
  }

  clear(): void {
    this.seenLinks.clear();
    this.saveCache();
    logger.info('Cache limpo completamente');
  }

  cleanOldEntries(daysOld: number = 30): void {
    logger.info('Limpeza de entradas antigas solicitada', { daysOld });
    logger.warn('Funcionalidade de limpeza por idade ainda não implementada');
  }

  getStats(): { totalLinks: number; cacheFile: string; exists: boolean } {
    return {
      totalLinks: this.seenLinks.size,
      cacheFile: this.cacheFile,
      exists: fs.existsSync(this.cacheFile)
    };
  }

  deleteCache(): void {
    try {
      if (fs.existsSync(this.cacheFile)) {
        fs.unlinkSync(this.cacheFile);
        logger.info('Arquivo de cache removido', { cacheFile: this.cacheFile });
      } else {
        logger.info('Arquivo de cache não existe', { cacheFile: this.cacheFile });
      }
      this.seenLinks.clear();
    } catch (error) {
      logger.error('Erro ao remover arquivo de cache', { 
        error: (error as Error).message,
        cacheFile: this.cacheFile 
      });
    }
  }
}