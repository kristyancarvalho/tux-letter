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

  private parseMarkdownToHtml(markdown: string): string {
    let html = markdown;
    
    html = html.replace(/^# (.*$)/gm, '<h1>$1</h1>');
    html = html.replace(/^## (.*$)/gm, '<h2>$1</h2>');
    html = html.replace(/^### (.*$)/gm, '<h3>$1</h3>');
    html = html.replace(/^#### (.*$)/gm, '<h4>$1</h4>');
    html = html.replace(/^\*\*(.+?)\*\*$/gm, '<strong>$1</strong>');
    html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
    html = html.replace(/\*(.+?)\*/g, '<em>$1</em>');
    html = html.replace(/`(.+?)`/g, '<code>$1</code>');
    html = html.replace(/^> (.*$)/gm, '<blockquote>$1</blockquote>');
    html = html.replace(/^\d+\.\s(.+$)/gm, '<li>$1</li>');
    html = html.replace(/^[\*\-\+]\s(.+$)/gm, '<li>$1</li>');
    html = html.replace(/\[(.+?)\]\((.+?)\)/g, '<a href="$2">$1</a>');
    html = html.replace(/```([^`]+)```/g, '<pre><code>$1</code></pre>');
    html = html.replace(/\n\n/g, '</p><p>');
    html = html.replace(/(\<li\>.*\<\/li\>)/gs, (match) => {
      if (match.includes('<li>') && !match.includes('<ol>') && !match.includes('<ul>')) {
        return '<ul>' + match + '</ul>';
      }
      return match;
    });
    
    if (!html.startsWith('<') && html.trim()) {
      html = '<p>' + html + '</p>';
    }
    
    return html;
  }

  private applyNewsletterStyling(html: string): string {
    return html
      .replace(/<h1>/g, '<h1 style="color: #1e293b; font-size: clamp(22px, 5vw, 28px); font-weight: 700; margin: 20px 0 16px 0; line-height: 1.3; border-bottom: 3px solid #3b82f6; padding-bottom: 10px;">')
      .replace(/<h2>/g, '<h2 style="color: #1e293b; font-size: clamp(20px, 4.5vw, 24px); font-weight: 600; margin: 20px 0 12px 0; line-height: 1.3; border-bottom: 2px solid #e5e7eb; padding-bottom: 8px;">')
      .replace(/<h3>/g, '<h3 style="color: #374151; font-size: clamp(18px, 4vw, 20px); font-weight: 600; margin: 16px 0 10px 0; line-height: 1.3;">')
      .replace(/<h4>/g, '<h4 style="color: #374151; font-size: clamp(16px, 3.5vw, 18px); font-weight: 600; margin: 14px 0 8px 0; line-height: 1.3;">')
      .replace(/<p>/g, '<p style="color: #374151; font-size: clamp(15px, 3.5vw, 16px); line-height: 1.6; margin: 0 0 14px 0; text-align: left;">')
      .replace(/<strong>/g, '<strong style="color: #1e293b; font-weight: 600;">')
      .replace(/<em>/g, '<em style="color: #4b5563; font-style: italic;">')
      .replace(/<code>/g, '<code style="background-color: #f3f4f6; color: #dc2626; padding: 3px 6px; border-radius: 4px; font-family: ui-monospace, SFMono-Regular, monospace; font-size: clamp(13px, 3vw, 14px); word-break: break-word;">')
      .replace(/<pre>/g, '<pre style="background-color: #1f2937; color: #f9fafb; padding: 15px; border-radius: 8px; overflow-x: auto; margin: 16px 0; border-left: 4px solid #3b82f6; -webkit-overflow-scrolling: touch;">')
      .replace(/<pre><code>/g, '<pre style="background-color: #1f2937; color: #f9fafb; padding: 15px; border-radius: 8px; overflow-x: auto; margin: 16px 0; border-left: 4px solid #3b82f6; -webkit-overflow-scrolling: touch;"><code style="background: none; color: inherit; padding: 0; font-family: ui-monospace, SFMono-Regular, monospace; font-size: clamp(13px, 3vw, 14px);">')
      .replace(/<ul>/g, '<ul style="color: #374151; font-size: clamp(15px, 3.5vw, 16px); line-height: 1.6; margin: 14px 0; padding-left: 18px;">')
      .replace(/<ol>/g, '<ol style="color: #374151; font-size: clamp(15px, 3.5vw, 16px); line-height: 1.6; margin: 14px 0; padding-left: 18px;">')
      .replace(/<li>/g, '<li style="margin-bottom: 6px;">')
      .replace(/<blockquote>/g, '<blockquote style="border-left: 4px solid #3b82f6; background-color: #f8fafc; padding: 12px 16px; margin: 16px 0; color: #4b5563; font-style: italic; border-radius: 0 8px 8px 0;">')
      .replace(/<a\s+href="([^"]*)"[^>]*>/g, '<a href="$1" style="color: #2563eb; text-decoration: none; font-weight: 500; word-break: break-word;">');
  }

  private formatEmailContent(data: EmailData): string {
    const currentDate = new Date().toLocaleDateString('pt-BR', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });

    const parsedContent = this.parseMarkdownToHtml(data.synthesizedText);
    const styledContent = this.applyNewsletterStyling(parsedContent);

    const referencesHtml = data.references.length > 0 
      ? data.references.map((link, index) => 
          `<tr>
            <td style="padding: 12px 0; border-bottom: 1px solid #f1f5f9;">
              <div style="display: flex; align-items: flex-start; gap: 12px;">
                <div style="color: white; font-size: 12px; font-weight: 600; padding: 4px; min-width: 24px; text-align: center; flex-shrink: 0; display: flex; align-items: center; justify-content: center;">${index + 1}</div>
                <a href="${link}" style="color: #2563eb; text-decoration: none; font-size: clamp(13px, 3vw, 14px); line-height: 1.4; word-break: break-all; flex: 1; display: block; padding-top: 2px;">${link}</a>
              </div>
            </td>
          </tr>`
        ).join('')
      : '<tr><td style="padding: 12px 0; color: #888; font-style: italic; text-align: center; font-size: clamp(13px, 3vw, 14px);">Nenhuma referência disponível</td></tr>';

    return `
      <!DOCTYPE html>
      <html lang="pt-BR">
      <head>
          <meta charset="UTF-8">
          <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes">
          <meta name="format-detection" content="telephone=no">
          <title>Tux Letter</title>
          <style>
            @media screen and (max-width: 480px) {
              .container { width: 100% !important; margin: 0 !important; }
              .header { padding: 30px 20px !important; }
              .header-title-container { flex-direction: column !important; gap: 10px !important; }
              .main-content { padding: 30px 20px !important; }
              .content-box { padding: 25px 20px !important; }
              .footer { padding: 20px !important; }
              .footer-content { flex-direction: column !important; gap: 15px !important; align-items: stretch !important; }
              .footer-stats { 
                display: grid !important; 
                grid-template-columns: repeat(2, 1fr) !important; 
                gap: 8px !important; 
                order: 1 !important;
              }
              .footer-date { 
                order: 2 !important; 
                align-self: center !important;
                width: fit-content !important;
              }
              .header-title { font-size: 32px !important; }
              .header-subtitle { font-size: 16px !important; }
              .references-header { padding: 15px 20px !important; }
              .references-content { padding: 20px !important; }
            }
            
            @media screen and (max-width: 360px) {
              .header { padding: 25px 15px !important; }
              .main-content { padding: 25px 15px !important; }
              .content-box { padding: 20px 15px !important; }
              .footer { padding: 15px !important; }
              .header-title { font-size: 28px !important; }
              .references-header { padding: 12px 15px !important; }
              .references-content { padding: 15px !important; }
              .footer-stats { 
                grid-template-columns: 1fr 1fr !important; 
                gap: 6px !important; 
              }
            }
          </style>
      </head>
      <body style="margin: 0; padding: 0; background-color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale;">
          <div class="container" style="max-width: 900px; margin: 0 auto; background-color: #ffffff; box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);">
              
              <header class="header" style="background: linear-gradient(135deg, #1e293b 0%, #334155 50%, #475569 100%); color: white; padding: 40px 30px; text-align: center; position: relative; overflow: hidden;">
                  <div style="position: absolute; top: 0; left: 0; right: 0; bottom: 0; opacity: 0.1; background-image: repeating-linear-gradient(45deg, transparent, transparent 10px, rgba(255,255,255,.05) 10px, rgba(255,255,255,.05) 20px);"></div>
                  <div style="position: relative; z-index: 1;">
                      <div style="display: inline-flex; align-items: center; justify-content: center; margin-bottom: 15px;">
                          <div style="background-color: rgba(255, 255, 255, 0.1); padding: 12px; border-radius: 12px; margin-right: 15px;">
                              <div style="font-size: 24px;">🐧</div>
                          </div>
                          <h1 style="margin: 0; font-size: 32px; font-weight: 700; letter-spacing: -0.5px;">Tux Letter</h1>
                      </div>
                      <p class="header-subtitle" style="margin: 0; font-size: clamp(14px, 4vw, 18px); opacity: 0.9; font-weight: 300; letter-spacing: 0.3px;">${currentDate}</p>
                  </div>
              </header>

              <main class="main-content" style="padding: 40px 30px;">
                  <div class="content-box" style="background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%); border: 1px solid #e2e8f0; margin-bottom: 40px; position: relative; overflow: hidden;">
                      <div style="position: absolute; top: 0; left: 0; width: 100%; height: 4px; background: linear-gradient(90deg, #3b82f6, #8b5cf6, #06b6d4); border-radius: 12px 12px 0 0;"></div>
                      <div style="color: #1e293b; line-height: 1.6;">
                          ${styledContent}
                      </div>
                  </div>

                  <div style="margin-bottom: 30px;">
                      <div class="references-header" style="background: linear-gradient(135deg, #1e293b 0%, #334155 100%); color: white; padding: 18px 25px; margin-bottom: 0;">
                          <h2 style="color: white; font-size: clamp(18px, 4vw, 22px); font-weight: 700; margin: 0; display: flex; align-items: center; gap: 12px;">
                              <div style="background: rgba(255,255,255,0.15); padding: 8px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0;">
                                  <span style="font-size: clamp(14px, 4vw, 18px); line-height: 1;">📎</span>
                              </div>
                              <span style="flex: 1;">Referências e Links</span>
                          </h2>
                      </div>
                      <div class="references-content" style="background-color: #ffffff; border: 1px solid #e2e8f0; border-top: none; padding: 25px;">
                          <table style="width: 100%; border-collapse: collapse;">
                              ${referencesHtml}
                          </table>
                      </div>
                  </div>
              </main>

              <footer class="footer" style="background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%); border-top: 1px solid #e5e7eb; padding: 25px 30px;">
                  <div class="footer-content" style="display: flex; justify-content: space-between; gap: 20px;">
                      <div class="footer-stats" style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; font-size: clamp(11px, 2.5vw, 13px); color: #64748b; font-weight: 500; flex: 1; max-width: 600px;">
                          <div style="display: flex; align-items: center; gap: 8px; padding: 10px 12px;">
                              <span style="font-size: clamp(12px, 3vw, 16px); flex-shrink: 0;">🤖</span>
                              <span style="white-space: nowrap;">${data.botVerificationCount} verif.</span>
                          </div>
                          <div style="display: flex; align-items: center; gap: 8px; padding: 10px 12px;">
                              <span style="font-size: clamp(12px, 3vw, 16px); flex-shrink: 0;">📰</span>
                              <span style="white-space: nowrap;">${data.totalItems} itens</span>
                          </div>
                          <div style="display: flex; align-items: center; gap: 8px; padding: 10px 12px;">
                              <span style="font-size: clamp(12px, 3vw, 16px); flex-shrink: 0;">📧</span>
                              <span style="white-space: nowrap;">${data.loreItems} lore</span>
                          </div>
                          <div style="display: flex; align-items: center; gap: 8px; padding: 10px 12px;">
                              <span style="font-size: clamp(12px, 3vw, 16px); flex-shrink: 0;">🔥</span>
                              <span style="white-space: nowrap;">${data.phoronixItems} phoronix</span>
                          </div>
                          <div style="display: flex; align-items: center; gap: 8px; padding: 10px 12px;">
                              <span style="font-size: clamp(12px, 3vw, 16px); flex-shrink: 0;">🐧</span>
                              <span style="white-space: nowrap;">${data.linuxcomItems} linux</span>
                          </div>
                          <div style="display: flex; align-items: center; gap: 8px; padding: 10px 12px;">
                              <span style="font-size: clamp(12px, 3vw, 16px); flex-shrink: 0;">🔗</span>
                              <span style="white-space: nowrap;">${data.itsfossItems} itsfoss</span>
                          </div>
                      </div>  
                      <div class="footer-date" style="font-size: clamp(10px, 2.5vw, 12px); color: #94a3b8; padding: 12px 16px; flex-shrink: 0; align-self: flex-start;">
                          <div style="font-weight: 600; margin-bottom: 4px; text-align: center;">Gerado em</div>
                          <div style="text-align: center; white-space: nowrap;">${new Date().toLocaleString('pt-BR')}</div>
                      </div>
                  </div>
                  <div style="text-align: center; margin-top: 20px; padding-top: 15px; border-top: 1px solid #e5e7eb;">
                      <p style="margin: 0; font-size: clamp(10px, 2.5vw, 12px); color: #9ca3af; font-weight: 400;">
                          Kristyan Carvalho • Tux Letter • Newsletter automatizada sobre Linux e Open Source
                      </p>
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