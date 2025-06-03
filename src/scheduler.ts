import * as cron from 'node-cron';
import { main } from './main';
import { logger } from './utils/logger';

class NewsScheduler {
  private job: cron.ScheduledTask | null = null;
  private isRunning = false;

  start() {
    this.job = cron.schedule('0 8 * * *', async () => {
      if (this.isRunning) {
        logger.warn('⏳ Job anterior ainda em execução, pulando esta iteração');
        return;
      }

      this.isRunning = true;
      logger.info('🕐 Executando job diário de scraping');

      try {
        await main();
        logger.info('✅ Job diário concluído com sucesso');
      } catch (error) {
        logger.error('❌ Erro no job diário', {
          error: (error as Error).message,
          stack: (error as Error).stack
        });
      } finally {
        this.isRunning = false;
      }
    }, {
      timezone: 'America/Sao_Paulo'
    });

    this.job.start();
    logger.info('📅 Scheduler iniciado - execução diária às 08:00');
  }

  stop() {
    if (this.job) {
      this.job.stop();
      this.job = null;
      logger.info('🛑 Scheduler parado');
    }
  }

  getStatus() {
    return {
      isScheduled: this.job !== null,
      isRunning: this.isRunning
    };
  }
}

const scheduler = new NewsScheduler();

process.on('SIGTERM', () => {
  logger.info('🔄 Recebido SIGTERM, parando scheduler');
  scheduler.stop();
  process.exit(0);
});

process.on('SIGINT', () => {
  logger.info('🔄 Recebido SIGINT, parando scheduler');
  scheduler.stop();
  process.exit(0);
});

if (require.main === module) {
  scheduler.start();
  
  logger.info('🚀 Aplicação iniciada com scheduler');
  logger.info('Status do scheduler:', scheduler.getStatus());
}

export { NewsScheduler, scheduler };