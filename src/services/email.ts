import nodemailer from 'nodemailer';
import { logger } from '../utils/logger';

interface EmailData {
  synthesizedText: string;
  references: string[];
  botVerificationCount: number;
  totalItems: number;
  loreItems: number;
  phoronixItems: number;
}

export class EmailService {
  private transporter: nodemailer.Transporter;
  private fromEmail: string;
  private toEmail: string = 'kristyancarvalho@gmail.com';

  constructor() {
    const gmailPassword = process.env.GMAIL_APP_PASSWORD;
    this.fromEmail = process.env.GMAIL_USER || 'tuxteller@gmail.com';

    if (!gmailPassword) {
      throw new Error('GMAIL_APP_PASSWORD não encontrada nas variáveis de ambiente');
    }

    this.transporter = nodemailer.createTransport({
      service: 'gmail',
      auth: {
        user: this.fromEmail,
        pass: gmailPassword
      }
    });
  }

  private formatEmailContent(data: EmailData): string {
    const referencesSection = data.references.length > 0 
      ? data.references.map((link, index) => `${index + 1}. ${link}`).join('\n')
      : 'Nenhuma referência disponível';

    const footer = `
═══════════════════════════════════════════
📊 INFORMAÇÕES DO SISTEMA
═══════════════════════════════════════════
🤖 Verificações de bot: ${data.botVerificationCount}
📰 Total de itens: ${data.totalItems}
📧 Mensagens lore.kernel.org: ${data.loreItems}
📰 Notícias Phoronix: ${data.phoronixItems}
⏰ Executado em: ${new Date().toLocaleString('pt-BR')}
═══════════════════════════════════════════`;

    return `TUX TELLER

${data.synthesizedText}

[LINKS DE REFERÊNCIA]
${referencesSection}

${footer}`;
  }

  async sendNewsEmail(data: EmailData): Promise<void> {
    try {
      logger.info('Preparando envio de email', {
        to: this.toEmail,
        totalItems: data.totalItems,
        referencesCount: data.references.length
      });

      const mailOptions = {
        from: this.fromEmail,
        to: this.toEmail,
        subject: `Tux Teller - ${new Date().toLocaleDateString('pt-BR')} - ${data.totalItems} notícias`,
        text: this.formatEmailContent(data)
      };

      const info = await this.transporter.sendMail(mailOptions);
      
      logger.info('Email enviado com sucesso', {
        messageId: info.messageId,
        to: this.toEmail,
        totalItems: data.totalItems
      });

    } catch (error) {
      logger.error('Erro ao enviar email', {
        error: (error as Error).message,
        to: this.toEmail
      });
      throw error;
    }
  }

  async testConnection(): Promise<boolean> {
    try {
      logger.info('Testando conexão SMTP');
      await this.transporter.verify();
      logger.info('Conexão SMTP verificada com sucesso');
      return true;
    } catch (error) {
      logger.error('Falha na verificação SMTP', {
        error: (error as Error).message
      });
      return false;
    }
  }
}