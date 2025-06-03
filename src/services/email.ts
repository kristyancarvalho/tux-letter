import nodemailer from 'nodemailer';
import dotenv from 'dotenv';
import path from 'path';
import { logger } from '../utils/logger';
import { EmailData } from '../types';

dotenv.config({ path: path.join(__dirname, '../../.env') });

export class EmailService {
  private transporter: nodemailer.Transporter;
  private fromEmail: string;
  private toEmail: string = 'kristyancarvalho@gmail.com';

  constructor() {
    const gmailPassword = process.env.GMAIL_APP_PASSWORD;
    this.fromEmail = process.env.GMAIL_USER || 'tuxletter@gmail.com';

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
    const currentDate = new Date().toLocaleDateString('pt-BR', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });

    const referencesHtml = data.references.length > 0 
      ? data.references.map((link, index) => 
          `<tr>
            <td style="padding: 8px 0; color: #555;">
              <strong style="color: #2563eb;">${index + 1}.</strong> 
              <a href="${link}" style="color: #2563eb; text-decoration: none; word-break: break-all;">${link}</a>
            </td>
          </tr>`
        ).join('')
      : '<tr><td style="padding: 8px 0; color: #888; font-style: italic;">Nenhuma referência disponível</td></tr>';

    return `
      <!DOCTYPE html>
      <html lang="pt-BR">
      <head>
          <meta charset="UTF-8">
          <meta name="viewport" content="width=device-width, initial-scale=1.0">
          <title>Tux Letter</title>
      </head>
      <body style="margin: 0; padding: 0; background-color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;">
          <div style="max-width: 800px; margin: 0 auto; background-color: #ffffff; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
              
              <header style="background: linear-gradient(135deg, #1e293b 0%, #334155 100%); color: white; padding: 40px 30px; text-align: center;">
                  <div style="display: inline-flex; align-items: center; justify-content: center; margin-bottom: 15px;">
                      <div style="background-color: rgba(255, 255, 255, 0.1); padding: 12px; border-radius: 12px; margin-right: 15px;">
                          <div style="font-size: 24px;">🐧</div>
                      </div>
                      <h1 style="margin: 0; font-size: 32px; font-weight: 700; letter-spacing: -0.5px;">Tux Letter</h1>
                  </div>
                  <p style="margin: 0; font-size: 16px; opacity: 0.9; font-weight: 300;">${currentDate}</p>
              </header>

              <main style="padding: 40px 30px;">
                  <div style="background-color: #f1f5f9; border-left: 4px solid #3b82f6; padding: 25px; margin-bottom: 40px; border-radius: 0 8px 8px 0;">
                      <div style="color: #1e293b; font-size: 16px; line-height: 1.7; white-space: pre-line;">${data.synthesizedText}</div>
                  </div>

                  <div style="margin-bottom: 30px;">
                      <h2 style="color: #1e293b; font-size: 20px; font-weight: 600; margin-bottom: 20px; padding-bottom: 10px; border-bottom: 2px solid #e5e7eb;">
                          📎 Referências
                      </h2>
                      <table style="width: 100%; border-collapse: collapse;">
                          ${referencesHtml}
                      </table>
                  </div>
              </main>

              <footer style="background-color: #f8fafc; border-top: 1px solid #e5e7eb; padding: 20px 30px;">
                  <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 15px;">
                      <div style="display: flex; flex-direction: column; gap: 5px; font-size: 12px; color: #64748b;">
                        <span>🤖 ${data.botVerificationCount} verificações</span>
                        <span>📰 ${data.totalItems} itens</span>
                        <span>📧 ${data.loreItems} lore</span>
                        <span>📰 ${data.phoronixItems} phoronix</span>
                        <span>🐧 ${data.linuxcomItems} linux.com</span>
                        <span>🔗 ${data.itsfossItems} itsfoss</span>
                      </div>  
                      <div style="font-size: 11px; color: #94a3b8;">
                          ${new Date().toLocaleString('pt-BR')}
                      </div>
                  </div>
              </footer>
          </div>
      </body>
      </html>`;
  }

  async sendNewsEmail(data: EmailData): Promise<void> {
    try {
      const currentDate = new Date().toLocaleDateString('pt-BR');
      
      logger.info('Preparando envio de email', {
        to: this.toEmail,
        totalItems: data.totalItems,
        referencesCount: data.references.length
      });

      const mailOptions = {
        from: `"Tux Letter" <${this.fromEmail}>`,
        to: this.toEmail,
        subject: `Tux Letter • ${currentDate} • ${data.totalItems} atualizações`,
        html: this.formatEmailContent(data)
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