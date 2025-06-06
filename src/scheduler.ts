import * as cron from 'node-cron';
import { main } from './main';
import { logger } from './utils/logger';

class NewsScheduler {
  private job: cron.ScheduledTask | null = null;
  private isRunning = false;

  start() {
    this.job = cron.schedule('0 20 * * *', async () => {
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
    logger.info('📅 Scheduler iniciado - execução diária às 20:00 (horário de São Paulo)');
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
      isRunning: this.isRunning,
      nextExecution: this.job ? '20:00 (todos os dias)' : 'N/A'
    };
  }
  async runNow() {
    if (this.isRunning) {
      logger.warn('⏳ Job já em execução, aguarde a conclusão');
      return false;
    }

    logger.info('🔄 Executando job manualmente');
    this.isRunning = true;

    try {
      await main();
      logger.info('✅ Job manual concluído com sucesso');
      return true;
    } catch (error) {
      logger.error('❌ Erro no job manual', {
        error: (error as Error).message,
        stack: (error as Error).stack
      });
      return false;
    } finally {
      this.isRunning = false;
    }
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
  
  logger.info('🚀 Scheduler iniciado com sucesso');
  logger.info('📊 Status do scheduler:', scheduler.getStatus());
  
  process.on('exit', () => {
    logger.info('👋 Encerrando scheduler');
  });
  
  setInterval(() => {
    logger.info('💓 Scheduler ativo - Status:', scheduler.getStatus());
  }, 60 * 60 * 1000);
}

export { NewsScheduler, scheduler };