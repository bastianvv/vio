type Listener = () => void;

const listeners = new Set<Listener>();

export function notifyScanFinished() {
  listeners.forEach((l) => l());
}

export function onScanFinished(listener: Listener) {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}
