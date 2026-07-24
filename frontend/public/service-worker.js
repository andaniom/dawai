const CACHE_NAME = 'dawai-v1';
const STATIC_ASSETS = [
  '/',
  '/index.html',
  '/_next/static',
];

// Install event: cache static assets
self.addEventListener('install', (event) => {
  console.log('[SW] Install event');
  self.skipWaiting();
});

// Activate event: cleanup old caches
self.addEventListener('activate', (event) => {
  console.log('[SW] Activate event');
  self.clients.claim();
  event.waitUntil(
    caches.keys().then((cacheNames) =>
      Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME)
          .map((name) => caches.delete(name))
      )
    )
  );
});

// Fetch event: cache-first for static, network-first for API
self.addEventListener('fetch', (event) => {
  const { request } = event;

  // Skip non-GET requests
  if (request.method !== 'GET') {
    return;
  }

  const url = new URL(request.url);

  // API calls: network-first (try network, fall back to cache)
  if (url.pathname.startsWith('/api/')) {
    return event.respondWith(
      fetch(request)
        .then((response) => {
          if (response.ok) {
            const cache = caches.open(CACHE_NAME);
            cache.then((c) => c.put(request, response.clone()));
          }
          return response;
        })
        .catch(() => {
          return caches.match(request);
        })
    );
  }

  // Static assets: cache-first (use cache, fall back to network)
  event.respondWith(
    caches.match(request).then((cached) => {
      return (
        cached ||
        fetch(request).then((response) => {
          if (response.ok && !url.pathname.includes('_next/data')) {
            caches.open(CACHE_NAME).then((cache) => {
              cache.put(request, response.clone());
            });
          }
          return response;
        })
      );
    })
  );
});
