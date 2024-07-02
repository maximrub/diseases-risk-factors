
export function truncateWords(string: string, words: number): string {
  const wordsArray = string.split(" ");
  if (wordsArray.length <= words) {
    return string;
  }
  return string.split(" ").splice(0, words).join(" ") + "...";
}

export function clamp(value: number, min: number, max: number): number {
  return Math.min(Math.max(value, min), max);
}