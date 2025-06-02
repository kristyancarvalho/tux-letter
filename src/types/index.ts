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