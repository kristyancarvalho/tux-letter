import { scrapeLoreMessages, scrapePhoronixNews, getCacheStats, persistCache } from './services/scraping';
import { OpenRouterService } from './services/openrouter';
import { logger } from './utils/logger';
import { Item } from './types';

async function main() {
  logger.info('🚀 Iniciando scraping de notícias Linux');

  const cacheStats = getCacheStats();
  logger.info('📊 Estatísticas do cache:', cacheStats);

  try {
    const openRouterService = new OpenRouterService();
    
    const connectionTest = await openRouterService.testConnection();
    if (!connectionTest) {
      logger.error('❌ Falha na conexão com OpenRouter API');
      throw new Error('OpenRouter API não está acessível');
    }
    logger.info('✅ Conexão com OpenRouter API confirmada');

    const loreItems = await scrapeLoreMessages();
    const phoronixItems = await scrapePhoronixNews();
    const allItems = [...loreItems, ...phoronixItems];

    const uniqueItems: Item[] = [];
    const seenLinks = new Set<string>();

    allItems.forEach(item => {
      if (!seenLinks.has(item.link)) {
        seenLinks.add(item.link);
        uniqueItems.push(item);
      }
    });

    if (uniqueItems.length === 0) {
      logger.info('📰 Nenhuma notícia nova encontrada');
      return {
        synthesizedText: 'Nenhuma notícia nova encontrada.',
        references: []
      };
    }

    logger.info('📰 Notícias novas encontradas:', { total: uniqueItems.length });
    
    uniqueItems.forEach((item, index) => {
      logger.info(`📄 Item coletado [${index + 1}/${uniqueItems.length}]`, {
        type: item.type,
        title: item.title,
        author: item.author,
        date: item.date,
        link: item.link,
        bodyPreview: item.body.substring(0, 100) + (item.body.length > 100 ? '...' : '')
      });
    });

    logger.info('🤖 Iniciando síntese de todas as notícias com OpenRouter API');
    const synthesizedNews = await openRouterService.synthesizeAllNews(uniqueItems);

    persistCache();

    const finalCacheStats = getCacheStats();
    logger.info('✅ Scraping e síntese concluídos', { 
      novasNoticias: uniqueItems.length,
      totalLinksCache: finalCacheStats.totalLinks,
      loreItems: loreItems.length,
      phoronixItems: phoronixItems.length,
      synthesizedTextLength: synthesizedNews.synthesizedText.length,
      totalReferences: synthesizedNews.references.length
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
  main().catch(error => {
    logger.error('❌ Erro no script principal', { 
      error: error.message,
      stack: error.stack 
    });
    process.exit(1);
  });
}

export { main };