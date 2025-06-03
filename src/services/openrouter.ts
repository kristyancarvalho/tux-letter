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
          max_tokens: 4000,
          temperature: 0.3
        },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json',
            'HTTP-Referer': 'https://localhost:3000',
            'X-Title': 'Tux Teller'
          },
          timeout: 45000
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
CONTEÚDO: ${item.body.substring(0, 2000)}
LINK: ${item.link}
---
`;
      });
    }

    if (newsItems.length > 0) {
      newsContent += '\n=== NOTÍCIAS GERAIS ===\n';
      newsItems.forEach((item, index) => {
        const source = this.getSourceFromLink(item.link);
        newsContent += `
NOTÍCIA ${index + 1} (${source}):
TÍTULO: ${item.title}
AUTOR: ${item.author}
DATA: ${item.date}
CONTEÚDO: ${item.body.substring(0, 2000)}
LINK: ${item.link}
---
`;
      });
    }

    const prompt = `
Você é um especialista em Linux e tecnologias open source com profundo conhecimento técnico. Analise todas as notícias e mensagens abaixo e crie um texto jornalístico abrangente em português brasileiro.

${newsContent}

INSTRUÇÕES PARA O TEXTO:

1. ESTRUTURA E ORGANIZAÇÃO:
   - Crie um texto corrido e coeso, como um artigo jornalístico especializado
   - Organize por temas e relevância (kernel, distribuições, aplicações, hardware, etc.)
   - Use transições naturais entre os assuntos
   - Mantenha fluidez narrativa sem divisões rígidas

2. CONTEÚDO E ANÁLISE:
   - Traduza e explique conceitos técnicos de forma acessível mas precisa
   - Contextualize as notícias dentro do ecossistema Linux/open source
   - Explique a importância e impacto de patches, atualizações e desenvolvimentos
   - Relacione diferentes notícias quando houver conexões temáticas
   - Inclua detalhes técnicos relevantes sem perder a clareza

3. ESTILO E TOM:
   - Use tom jornalístico informativo e técnico, mas acessível
   - Seja objetivo e factual, evitando especulações
   - Mantenha interesse do leitor com linguagem envolvente
   - Use terminologia técnica correta em português e inglês quando necessário

4. EXTENSÃO E PROFUNDIDADE:
   - Crie um texto substancial de pelo menos 800-1200 palavras
   - Desenvolva cada tópico com profundidade adequada
   - Não seja superficial - explore as implicações das notícias
   - Inclua contexto histórico quando relevante

5. INTEGRAÇÃO DE FONTES:
   - Integre naturalmente informações de lore.kernel.org, Phoronix, Linux.com e Its FOSS
   - Mencione a fonte apenas quando necessário para credibilidade
   - Trate patches do kernel com especial atenção técnica
   - Balance notícias de diferentes fontes harmoniosamente

IMPORTANTE: Certifique-se de abordar TODAS as notícias fornecidas, integrando-as em um texto coeso e informativo que sirva como um resumo completo das novidades mais relevantes do mundo Linux.

Responda APENAS com o texto sintetizado, sem preâmbulos, cabeçalhos ou divisões especiais.
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

  private getSourceFromLink(link: string): string {
    if (link.includes('phoronix.com')) return 'Phoronix';
    if (link.includes('linux.com')) return 'Linux.com';
    if (link.includes('itsfoss.com')) return 'Its FOSS';
    if (link.includes('lore.kernel.org')) return 'lore.kernel.org';
    return 'Fonte desconhecida';
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