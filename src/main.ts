import { ScraperManager } from './scrapers';
import { OpenRouterService } from './services/openrouter';
import { EmailService } from './services/email';
import { logger } from './utils/logger';
import { scheduler } from './scheduler';
import * as http from 'http';

async function main() {
  logger.info('🚀 Iniciando scraping de notícias Linux');

  const scraperManager = new ScraperManager();
  const cacheStats = scraperManager.getCacheStats();
  logger.info('📊 Estatísticas do cache:', cacheStats);

  try {
    const openRouterService = new OpenRouterService();
    const emailService = new EmailService();
    
    const connectionTest = await openRouterService.testConnection();
    if (!connectionTest) {
      logger.error('❌ Falha na conexão com OpenRouter API');
      throw new Error('OpenRouter API não está acessível');
    }
    logger.info('✅ Conexão com OpenRouter API confirmada');

    const emailTest = await emailService.testConnection();
    if (!emailTest) {
      logger.error('❌ Falha na conexão SMTP');
      throw new Error('SMTP não está acessível');
    }
    logger.info('✅ Conexão SMTP confirmada');

    const allItems = await scraperManager.scrapeAll();

    if (allItems.length === 0) {
      logger.info('📰 Nenhuma notícia nova encontrada');
      
      await emailService.sendNewsEmail({
        synthesizedText: 'Nenhuma notícia nova encontrada hoje.',
        references: [],
        botVerificationCount: scraperManager.getBotVerificationCount(),
        totalItems: 0,
        loreItems: 0,
        phoronixItems: 0,
        linuxcomItems: 0,
        itsfossItems: 0
      });

      return {
        synthesizedText: 'Nenhuma notícia nova encontrada.',
        references: []
      };
    }

    logger.info('📰 Notícias novas encontradas:', { total: allItems.length });
    
    allItems.forEach((item, index) => {
      logger.info(`📄 Item coletado [${index + 1}/${allItems.length}]`, {
        type: item.type,
        title: item.title,
        author: item.author,
        date: item.date,
        link: item.link,
        bodyPreview: item.body.substring(0, 100) + (item.body.length > 100 ? '...' : '')
      });
    });

    logger.info('🤖 Iniciando síntese de todas as notícias com OpenRouter API');
    const synthesizedNews = await openRouterService.synthesizeAllNews(allItems);

    const loreItems = allItems.filter(item => item.type === 'patch' || item.type === 'inbox');
    const phoronixItems = allItems.filter(item => item.link.includes('phoronix.com'));
    const linuxcomItems = allItems.filter(item => item.link.includes('linux.com'));
    const itsfossItems = allItems.filter(item => item.link.includes('itsfoss.com'));

    logger.info('📧 Enviando email com notícias sintetizadas');
    await emailService.sendNewsEmail({
      synthesizedText: synthesizedNews.synthesizedText,
      references: synthesizedNews.references,
      botVerificationCount: scraperManager.getBotVerificationCount(),
      totalItems: allItems.length,
      loreItems: loreItems.length,
      phoronixItems: phoronixItems.length,
      linuxcomItems: linuxcomItems.length,
      itsfossItems: itsfossItems.length
    });

    scraperManager.persistCache();

    const finalCacheStats = scraperManager.getCacheStats();
    logger.info('✅ Scraping, síntese e envio de email concluídos', { 
      novasNoticias: allItems.length,
      totalLinksCache: finalCacheStats.totalLinks,
      loreItems: loreItems.length,
      phoronixItems: phoronixItems.length,
      linuxcomItems: linuxcomItems.length,
      itsfossItems: itsfossItems.length,
      synthesizedTextLength: synthesizedNews.synthesizedText.length,
      totalReferences: synthesizedNews.references.length,
      emailSent: true
    });

    return synthesizedNews;

  } catch (error) {
    logger.error('❌ Erro durante o processamento principal', {
      error: (error as Error).message,
      stack: (error as Error).stack
    });
    throw error;
  }
}

if (require.main === module) {
  const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ 
      status: 'running',
      scheduler: scheduler.getStatus(),
      timestamp: new Date().toISOString()
    }));
  });

  const PORT = process.env.PORT || 8080;
  server.listen(PORT, () => {
    logger.info(`🌐 Fake HTTP Server iniciado na porta ${PORT}`);
  });

  scheduler.start();
  logger.info('🚀 Aplicação iniciada com scheduler');
  logger.info('Status do scheduler:', scheduler.getStatus());
}

export { main };