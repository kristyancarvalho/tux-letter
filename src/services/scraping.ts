import { ScraperManager } from '../scrapers';
import { Item } from '../types';

const scraperManager = new ScraperManager();

export async function scrapeLoreMessages(): Promise<Item[]> {
  return await scraperManager.scrapeSpecific(['lore']);
}

export async function scrapePhoronixNews(): Promise<Item[]> {
  return await scraperManager.scrapeSpecific(['phoronix']);
}

export async function scrapeLinuxComNews(): Promise<Item[]> {
  return await scraperManager.scrapeSpecific(['linuxcom']);
}

export async function scrapeItsFossNews(): Promise<Item[]> {
  return await scraperManager.scrapeSpecific(['itsfoss']);
}

export async function scrapeAllNews(): Promise<Item[]> {
  return await scraperManager.scrapeAll();
}

export function getCacheStats() {
  return scraperManager.getCacheStats();
}

export function clearCache() {
  scraperManager.clearCache();
}

export function persistCache() {
  scraperManager.persistCache();
}

export function getBotVerificationCount(): number {
  return scraperManager.getBotVerificationCount();
}