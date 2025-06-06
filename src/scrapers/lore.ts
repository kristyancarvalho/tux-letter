import * as cheerio from 'cheerio';
import { logger } from '../utils/logger';
import { Item } from '../types';
import { BaseScraper } from './base';
import { NewsCache } from '../utils/cache';

export class LoreScraper extends BaseScraper {
  private baseUrl = 'https://lore.kernel.org';
  private botVerificationCount = 0;

  constructor(sharedCache?: NewsCache) {
    super(sharedCache);
  }

  async scrape(): Promise<Item[]> {
    const url = `${this.baseUrl}/lkml/`;
    logger.info('Iniciando scraping do lore.kernel.org', { url });

    try {
      const html = await this.fetchPage(url);
      
      if (html.includes("Making sure you're not a bot!")) {
        this.botVerificationCount++;
        logger.warn('Verificação de bot detectada no lore.kernel.org', { 
          count: this.botVerificationCount 
        });
        return [];
      }

      const $ = cheerio.load(html);
      const items: Item[] = [];

      const messageLinks = $('tbody tr').map((i, row) => {
        const $row = $(row);
        const $link = $row.find('td:first-child a');
        const href = $link.attr('href');
        const title = $link.text().trim();
        
        if (href && title && !href.includes('archive')) {
          return {
            href: href.startsWith('http') ? href : `${this.baseUrl}${href}`,
            title: title
          };
        }
        return null;
      }).get().filter(Boolean);

      if (messageLinks.length === 0) {
        logger.warn('Nenhum link de mensagem encontrado no lore.kernel.org');
        return [];
      }

      const hrefs = messageLinks.map(link => link.href);
      const newLinks = this.getUniqueNewLinks(hrefs);
      
      this.logLinkStats('lore.kernel.org', messageLinks.length, hrefs.length, newLinks.length);

      if (newLinks.length === 0) {
        return [];
      }

      for (const href of newLinks.slice(0, this.maxItems)) {
        const linkData = messageLinks.find(link => link.href === href);
        if (!linkData) continue;

        if (!this.cache.isNew(href)) continue;

        this.cache.markAsSeen(href);

        const item: Item = {
          type: linkData.title.includes('[PATCH]') ? 'patch' : 'inbox',
          title: linkData.title,
          author: 'Desconhecido',
          date: 'Desconhecida',
          body: '',
          link: href
        };

        try {
          await this.delay(1500);
          const messageHtml = await this.fetchPage(href);
          const $msg = cheerio.load(messageHtml);

          const subject = $msg('span:contains("Subject:")').parent().text().replace('Subject:', '').trim();
          const from = $msg('span:contains("From:")').parent().text().replace('From:', '').trim();
          const date = $msg('span:contains("Date:")').parent().text().replace('Date:', '').trim();
          const bodyText = $msg('pre').first().text().trim();

          item.title = subject || item.title;
          item.author = from || item.author;
          item.date = date || item.date;
          item.body = bodyText || 'Conteúdo não disponível';

          items.push(item);

          logger.info('Mensagem processada do lore.kernel.org', {
            title: item.title.substring(0, 50),
            author: item.author,
            type: item.type
          });

        } catch (error) {
          logger.error('Erro ao processar mensagem do lore.kernel.org', {
            url: href,
            error: (error as Error).message
          });
        }
      }

      logger.info('Scraping do lore.kernel.org concluído', { 
        itemsColetados: items.length 
      });

      return items;

    } catch (error) {
      logger.error('Erro no scraping do lore.kernel.org', {
        error: (error as Error).message
      });
      return [];
    }
  }

  getBotVerificationCount(): number {
    return this.botVerificationCount;
  }
}