import axios from 'axios';
import { logger } from '../utils/logger';
import { Item } from '../types';
import dotenv from 'dotenv';
import path from 'path';

dotenv.config({ path: path.join(__dirname, '../../.env') });

interface OpenRouterResponse {
  choices: Array<{
    message: {
      content: string;
    };
  }>;
}

interface SynthesizedNews {
  synthesizedText: string;
  references: string[];
}

export class OpenRouterService {
  private apiKey: string;
  private baseURL = 'https://openrouter.ai/api/v1/chat/completions';
  private model = 'meta-llama/llama-3.1-8b-instruct:free';

  constructor() {
    this.apiKey = process.env.OPENROUTER_API_KEY || '';
    if (!this.apiKey) {
      throw new Error('OPENROUTER_API_KEY não encontrada nas variáveis de ambiente');
    }
  }

  private async makeRequest(prompt: string): Promise<string> {
    try {
      logger.info('Fazendo requisição para OpenRouter', {
        model: this.model,
        promptLength: prompt.length
      });

      const response = await axios.post<OpenRouterResponse>(
        this.baseURL,
        {
          model: this.model,
          messages: [
            {
              role: 'user',
              content: prompt
            }
          ],
          max_tokens: 2000,
          temperature: 0.3
        },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json',
            'HTTP-Referer': 'https://localhost:3000',
            'X-Title': 'Tux Teller'
          },
          timeout: 30000
        }
      );

      logger.info('Resposta recebida do OpenRouter', {
        status: response.status,
        responseLength: response.data.choices[0]?.message?.content?.length || 0
      });

      return response.data.choices[0]?.message?.content || '';

    } catch (error) {
      if (axios.isAxiosError(error)) {
        logger.error('Erro na requisição para OpenRouter', {
          status: error.response?.status,
          statusText: error.response?.statusText,
          data: error.response?.data,
          message: error.message
        });
      } else {
        logger.error('Erro desconhecido na requisição para OpenRouter', {
          error: (error as Error).message
        });
      }
      throw error;
    }
  }

  async synthesizeAllNews(items: Item[]): Promise<SynthesizedNews> {
    if (items.length === 0) {
      return {
        synthesizedText: 'Nenhuma notícia nova encontrada.',
        references: []
      };
    }

    const loreItems = items.filter(item => item.type === 'patch' || item.type === 'inbox');
    const newsItems = items.filter(item => item.type === 'news');

    let newsContent = '';
    
    if (loreItems.length > 0) {
      newsContent += '\n=== MENSAGENS DO KERNEL LINUX (lore.kernel.org) ===\n';
      loreItems.forEach((item, index) => {
        newsContent += `
MENSAGEM ${index + 1}:
TIPO: ${item.type}
TÍTULO: ${item.title}
AUTOR: ${item.author}
DATA: ${item.date}
CONTEÚDO: ${item.body.substring(0, 1500)}
LINK: ${item.link}
---
`;
      });
    }

    if (newsItems.length > 0) {
      newsContent += '\n=== NOTÍCIAS GERAIS (Phoronix) ===\n';
      newsItems.forEach((item, index) => {
        newsContent += `
NOTÍCIA ${index + 1}:
TIPO: ${item.type}
TÍTULO: ${item.title}
AUTOR: ${item.author}
DATA: ${item.date}
CONTEÚDO: ${item.body.substring(0, 1500)}
LINK: ${item.link}
---
`;
      });
    }

    const prompt = `
Você é um especialista em Linux e tecnologia. Analise todas as notícias abaixo e crie um texto sintetizado em português brasileiro.

${newsContent}

Tarefas:
1. Crie um texto corrido em português brasileiro que sintetize todas as notícias
2. Organize por temas quando possível (patches do kernel, notícias gerais, etc.)
3. Traduza e resuma o conteúdo técnico de forma clara
4. Mantenha o foco nos aspectos mais relevantes para a comunidade Linux
5. Use um tom informativo e técnico apropriado
6. IMPORTANTE: Certifique-se de incluir informações tanto das mensagens do lore.kernel.org quanto das notícias do Phoronix

Formato de resposta:
Apenas o texto sintetizado, sem cabeçalhos ou divisões especiais.
`;

    try {
      logger.info('Iniciando síntese de todas as notícias', {
        totalItems: items.length,
        loreItems: loreItems.length,
        newsItems: newsItems.length,
        promptLength: prompt.length
      });

      const response = await this.makeRequest(prompt);
      const references = items.map(item => item.link);

      logger.info('📰 Resposta da API:');
      logger.info('[TEXTO DA LLM]');
      logger.info(response);
      logger.info('[LINKS DAS REFERENCIAS]');
      references.forEach((link, index) => {
        logger.info(`${index + 1}. ${link}`);
      });

      logger.info('Síntese de notícias concluída', {
        responseLength: response.length,
        totalReferences: references.length,
        loreReferences: loreItems.length,
        newsReferences: newsItems.length
      });

      return {
        synthesizedText: response,
        references: references
      };

    } catch (error) {
      logger.error('Erro ao sintetizar notícias com OpenRouter', {
        error: (error as Error).message,
        totalItems: items.length
      });
      
      return {
        synthesizedText: 'Não foi possível sintetizar as notícias devido a erro na API',
        references: items.map(item => item.link)
      };
    }
  }

  async testConnection(): Promise<boolean> {
    try {
      logger.info('Testando conexão com OpenRouter');
      
      const testPrompt = 'Responda apenas "OK" se você conseguir me ouvir.';
      const response = await this.makeRequest(testPrompt);
      
      const isWorking = response.toLowerCase().includes('ok');
      
      logger.info('Teste de conexão finalizado', {
        success: isWorking,
        response: response.substring(0, 100)
      });
      
      return isWorking;
      
    } catch (error) {
      logger.error('Falha no teste de conexão', {
        error: (error as Error).message
      });
      return false;
    }
  }
}