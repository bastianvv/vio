type EnrichHandler = () => Promise<void>;

let handler: EnrichHandler | null = null;

export function registerEnrichHandler(h: EnrichHandler | null) {
  handler = h;
}

export async function triggerEnrich() {
  if (!handler) return;
  await handler();
}

export function hasEnrichHandler() {
  return handler !== null;
}
