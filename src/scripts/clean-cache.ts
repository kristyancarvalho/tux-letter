import { NewsCache } from '../utils/cache';
import { logger } from '../utils/logger';

async function cleanCache() {
  logger.info('🧹 Iniciando limpeza do cache...');
  
  try {
    const cache = new NewsCache();
    
    const statsBefore = cache.getStats();
    logger.info('📊 Estatísticas antes da limpeza:', {
      totalLinks: statsBefore.totalLinks,
      cacheExists: statsBefore.exists,
      cacheFile: statsBefore.cacheFile
    });
    
    cache.deleteCache();
    
    const statsAfter = cache.getStats();
    logger.info('✅ Cache limpo com sucesso:', {
      totalLinksRemovidos: statsBefore.totalLinks,
      cacheExistsAfter: statsAfter.exists
    });
    
    logger.info('🎉 Limpeza concluída! O próximo scraping processará todas as notícias novamente.');
    
  } catch (error) {
    logger.error('❌ Erro durante a limpeza do cache:', {
      error: (error as Error).message,
      stack: (error as Error).stack
    });
    process.exit(1);
  }
}

if (require.main === module) {
  cleanCache()
    .then(() => process.exit(0))
    .catch((error) => {
      logger.error('Erro fatal no script de limpeza:', { error: error.message });
      process.exit(1);
    });
}

export { cleanCache };