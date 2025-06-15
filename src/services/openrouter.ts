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
            'X-Title': 'Tux Letter'
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
Você é um editor especializado em newsletters sobre Linux e tecnologias open source. Crie uma newsletter completa em português brasileiro usando Markdown, com base nas notícias e mensagens fornecidas abaixo.

${newsContent}

INSTRUÇÕES PARA A NEWSLETTER:

## ESTRUTURA EM MARKDOWN:

1. **CABEÇALHO PRINCIPAL:**
   - Use \`# 📰 Destaques das últimas 24 horas\` como título principal
   - Inclua uma breve introdução sobre o que será abordado

2. **SEÇÕES ORGANIZADAS:**
   - \`## 🔧 Desenvolvimento do Kernel Linux\` (para patches e lore)
   - \`## 📰 Notícias em Destaque\` (para notícias gerais)
   - \`## 🚀 Lançamentos e Atualizações\` (se houver)
   - \`## 💡 Tecnologias Emergentes\` (se relevante)

3. **FORMATAÇÃO MARKDOWN:**
   - Use **negrito** para destacar pontos importantes
   - Use *itálico* para ênfase
   - Use \`código\` para termos técnicos, versões e comandos
   - Use > citações para destacar informações importantes
   - Use listas numeradas ou com marcadores quando apropriado

## CONTEÚDO E ESTILO:

1. **Tom de Newsletter:**
   - Linguagem acessível mas tecnicamente precisa
   - Engajamento direto com o leitor
   - Contextualize desenvolvimentos dentro do ecossistema Linux
   - Explique a importância prática de cada notícia

2. **Organização do Conteúdo:**
   - Comece com um resumo executivo dos principais destaques
   - Agrupe notícias relacionadas em seções temáticas
   - Desenvolva cada tópico com profundidade adequada
   - Inclua detalhes técnicos relevantes sem perder clareza

3. **Análise e Contexto:**
   - Explique o impacto de patches e atualizações do kernel
   - Relacione diferentes notícias quando houver conexões
   - Inclua perspectivas sobre tendências e desenvolvimentos futuros
   - Mencione benefícios práticos para usuários e desenvolvedores

## DIRETRIZES ESPECÍFICAS:

- **Extensão:** 1000-1500 palavras, bem estruturadas
- **Patches do Kernel:** Explique funcionalidades, melhorias de performance e correções
- **Notícias Gerais:** Contextualize dentro do ecossistema open source
- **Linguagem Técnica:** Use terminologia correta em português, com termos em inglês quando necessário
- **Integração:** Todas as notícias devem ser abordadas de forma coesa

## EXEMPLO DE ESTRUTURA:

\`\`\`markdown
# 🐧 Destaques das últimas 24 horas

Bem-vindos à mais nova edição da Tux Letter! O último dia trouxe desenvolvimentos significativos...

## 🔧 Desenvolvimento do Kernel Linux

### Melhorias no Subsistema de Rede

O kernel Linux recebeu importantes atualizações...

## 📰 Notícias em Destaque

### Nova Versão do Ubuntu

A Canonical anunciou...

## 💡 Perspectivas

Com esses desenvolvimentos, podemos esperar...
\`\`\`

IMPORTANTE: 
- Responda APENAS com o conteúdo da newsletter em Markdown
- Não inclua preâmbulos ou explicações sobre o formato
- Assegure-se de abordar TODAS as notícias fornecidas
- Mantenha consistência no uso de emojis e formatação
- Crie um fluxo narrativo envolvente e informativo

Responda exclusivamente com a newsletter em formato Markdown.
`;

    try {
      logger.info('Iniciando síntese de newsletter', {
        totalItems: items.length,
        loreItems: loreItems.length,
        newsItems: newsItems.length,
        promptLength: prompt.length
      });

      const response = await this.makeRequest(prompt);
      const references = items.map(item => item.link);

      logger.info('📰 Resposta da API:');
      logger.info('[NEWSLETTER EM MARKDOWN]');
      logger.info(response);
      logger.info('[LINKS DAS REFERENCIAS]');
      references.forEach((link, index) => {
        logger.info(`${index + 1}. ${link}`);
      });

      logger.info('Newsletter gerada com sucesso', {
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
      logger.error('Erro ao gerar newsletter com OpenRouter', {
        error: (error as Error).message,
        totalItems: items.length
      });
      
      return {
        synthesizedText: `# 🐧 Tux Letter

## ⚠️ Erro na Geração

Não foi possível gerar a newsletter devido a um erro na API. Por favor, tente novamente mais tarde.

## 📊 Estatísticas

- **Total de itens:** ${items.length}
- **Mensagens do kernel:** ${loreItems.length}
- **Notícias gerais:** ${newsItems.length}`,
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
