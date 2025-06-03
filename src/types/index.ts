export interface Item {
  type: 'inbox' | 'patch' | 'news';
  title: string;
  author: string;
  date: string;
  body: string;
  link: string;
}

export interface SynthesizedNews {
  synthesizedText: string;
  references: string[];
}

export interface EmailData {
  synthesizedText: string;
  references: any[];
  botVerificationCount: number;
  totalItems: number;
  loreItems: number;
  phoronixItems: number;
  linuxcomItems: number;
  itsfossItems: number;
}

export interface BaseScraper {
  scrape(): Promise<Item[]>;
}