import axios from 'axios';
import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { NewsCache } from '../utils/cache';
import { Item } from '../types';

let botVerificationCount = 0;
const cache = new NewsCache();

export async function scrapeLoreMessages(): Promise<Item[]> {
  const url = 'https://lore.kernel.org/lkml/';
  logger.info('Acessando página principal do lore.kernel.org', { url });

  try {
    const response = await axios.get(url, {
      headers: {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Connection': 'keep-alive'
      }
    });
    logger.info('Página principal carregada', { status: response.status, length: response.data.length });

    if (response.data.includes("Making sure you're not a bot!")) {
      botVerificationCount++;
      logger.warn('Verificação de bot detectada no lore.kernel.org', { botVerificationCount });
      return [];
    }

    const $ = cheerio.load(response.data);
    const items: Item[] = [];

    const links = $('a[href*="@kernel.org"]');
    logger.info('Links encontrados antes do filtro no lore.kernel.org', { count: links.length });

    const filteredLinks = links.filter((i, el) => {
      const href = $(el).attr('href') || '';
      return !href.includes('archive');
    });

    logger.info('Links de mensagens encontrados no lore.kernel.org', { count: filteredLinks.length });

    const allHrefs = filteredLinks.map((i, el) => $(el).attr('href')).get().filter(href => href);
    const uniqueLinks = Array.from(new Set(allHrefs));
    
    const newLinks = cache.filterNewLinks(uniqueLinks).slice(0, 5);
    logger.info('Links novos após filtro de cache no lore.kernel.org', { 
      totalUnique: uniqueLinks.length,
      newLinks: newLinks.length,
      cached: uniqueLinks.length - newLinks.length
    });

    if (newLinks.length === 0) {
      logger.info('Nenhum link novo encontrado no lore.kernel.org');
      return [];
    }

    const processedLinks = new Set<string>();
    
    filteredLinks.each((i, el) => {
      const href = $(el).attr('href') || '';
      const fullUrl = href.startsWith('http') ? href : `https://lore.kernel.org${href}`;
      
      if (newLinks.includes(href) && items.length < 5 && !processedLinks.has(fullUrl)) {
        processedLinks.add(fullUrl);
        items.push({
          type: $(el).text().includes('[PATCH]') ? 'patch' : 'inbox',
          title: $(el).text().trim() || 'Sem título',
          author: 'Desconhecido',
          date: 'Desconhecida',
          body: '',
          link: fullUrl
        });
      }
    });

    logger.info('Items únicos coletados do lore.kernel.org', { count: items.length });

    for (const item of items) {
      try {
        logger.info('Acessando página da mensagem lore.kernel.org', { url: item.link });
        
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        const messageResponse = await axios.get(item.link, {
          headers: {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
          }
        });
        const $msg = cheerio.load(messageResponse.data);

        item.title = $msg('h1').text().trim() || item.title;
        item.author = $msg('b:contains("From:")').next().text().trim() || 'Desconhecido';
        item.date = $msg('b:contains("Date:")').next().text().trim() || 'Desconhecida';
        item.body = $msg('pre').text().trim() || 'Sem corpo';
        
        cache.markAsSeen(item.link);
        
      } catch (error) {
        logger.error('Erro ao acessar mensagem lore.kernel.org', {
          url: item.link,
          error: (error as Error).message,
          botVerificationCount
        });
      }
    }

    return items;
    
  } catch (error) {
    logger.error('Erro ao fazer scraping da página principal do lore.kernel.org', {
      url,
      error: (error as Error).message,
      botVerificationCount
    });
    return [];
  }
}

export async function scrapePhoronixNews(): Promise<Item[]> {
  const url = 'https://www.phoronix.com/';
  logger.info('Acessando página principal do Phoronix', { url });

  try {
    const response = await axios.get(url, {
      headers: {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Connection': 'keep-alive'
      }
    });
    logger.info('Página principal carregada', { status: response.status, length: response.data.length });

    const $ = cheerio.load(response.data);
    const items: Item[] = [];

    const links = $('a[href*="/news/"]');
    logger.info('Links encontrados antes do filtro no Phoronix', { count: links.length });

    const filteredLinks = links.filter((i, el) => {
      const href = $(el).attr('href') || '';
      return href.includes('/news/') && !href.includes('rss') && $(el).closest('.sidebar').length === 0;
    });

    logger.info('Links de notícias encontrados no Phoronix', { count: filteredLinks.length });

    const allHrefs = filteredLinks.map((i, el) => $(el).attr('href')).get().filter(href => href);
    const uniqueLinks = Array.from(new Set(allHrefs));
    
    const newLinks = cache.filterNewLinks(uniqueLinks).slice(0, 5);
    logger.info('Links novos após filtro de cache no Phoronix', { 
      totalUnique: uniqueLinks.length,
      newLinks: newLinks.length,
      cached: uniqueLinks.length - newLinks.length
    });

    if (newLinks.length === 0) {
      logger.info('Nenhum link novo encontrado no Phoronix');
      return [];
    }

    const processedLinks = new Set<string>();

    filteredLinks.each((i, el) => {
      const href = $(el).attr('href') || '';
      const fullUrl = href.startsWith('http') ? href : `https://www.phoronix.com${href}`;
      
      if (newLinks.includes(href) && items.length < 5 && !processedLinks.has(fullUrl)) {
        processedLinks.add(fullUrl);
        items.push({
          type: 'news',
          title: $(el).text().trim() || 'Sem título',
          author: 'Desconhecido',
          date: 'Desconhecida',
          body: '',
          link: fullUrl
        });
      }
    });

    logger.info('Items únicos coletados do Phoronix', { count: items.length });

    for (const item of items) {
      try {
        logger.info('Acessando página da notícia Phoronix', { url: item.link });
        
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        const newsResponse = await axios.get(item.link, {
          headers: {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
          }
        });
        const $news = cheerio.load(newsResponse.data);

        item.title = $news('h1').text().trim() || item.title;
        item.author = $news('.author, .byline, [class*="author"]').first().text().trim() || 'Desconhecido';
        item.date = $news('.date, .published, time, [datetime]').first().text().trim() || 'Desconhecida';
        item.body = $news('article, .content').text().trim() || 'Sem corpo';
        
        cache.markAsSeen(item.link);
        
      } catch (error) {
        logger.error('Erro ao acessar notícia Phoronix', {
          url: item.link,
          error: (error as Error).message
        });
      }
    }

    return items;
    
  } catch (error) {
    logger.error('Erro ao fazer scraping da página principal do Phoronix', {
      url,
      error: (error as Error).message
    });
    return [];
  }
}

export function getCacheStats() {
  return cache.getStats();
}

export function clearCache() {
  cache.clear();
}

export function persistCache() {
  cache.persist();
}