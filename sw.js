const cacheName = self.location.pathname
const pages = [

  "/propertieslint/rules/comment/comment-spaces/",
    "/propertieslint/rules/comment/comment-style/",
    "/propertieslint/rules/comment/",
    "/propertieslint/rules/format/invalid-escapes/",
    "/propertieslint/rules/format/missing-separator/",
    "/propertieslint/rules/format/unterminated-line-continuation/",
    "/propertieslint/rules/format/",
    "/propertieslint/rules/pair/duplicate-key/",
    "/propertieslint/rules/pair/key-name/",
    "/propertieslint/rules/pair/missing-key/",
    "/propertieslint/rules/pair/missing-value/",
    "/propertieslint/rules/pair/",
    "/propertieslint/rules/whitespace/duplicate-blank-line/",
    "/propertieslint/rules/whitespace/no-leading-blank-line/",
    "/propertieslint/rules/whitespace/trailing-newline/",
    "/propertieslint/rules/whitespace/untrimmed-entry/",
    "/propertieslint/rules/whitespace/",
    "/propertieslint/",
    "/propertieslint/book/",
    "/propertieslint/categories/",
    "/propertieslint/rules/",
    "/propertieslint/tags/",
    "/propertieslint/book.min.98c83e4c8e1c8661368bcb29bbaf4dcc9d31b8164f76aae322d82d86a9750c65.css",
  "/propertieslint/en.search-data.min.b858a452a6e60089c21b87c87d21e235365bcbfffcac6eea1ba25738a911a708.json",
  "/propertieslint/en.search.min.bdd1203f9e10005b6b2d6c8ff5eb48a149084cd6ba87553d074a3d79eabf0f2a.js",
  
];

self.addEventListener("install", function (event) {
  self.skipWaiting();

  caches.open(cacheName).then((cache) => {
    return cache.addAll(pages);
  });
});

self.addEventListener("fetch", (event) => {
  const request = event.request;
  if (request.method !== "GET") {
    return;
  }

  /**
   * @param {Response} response
   * @returns {Promise<Response>}
   */
  function saveToCache(response) {
    if (cacheable(response)) {
      return caches
        .open(cacheName)
        .then((cache) => cache.put(request, response.clone()))
        .then(() => response);
    } else {
      return response;
    }
  }

  /**
   * @param {Error} error
   */
  function serveFromCache(error) {
    return caches.open(cacheName).then((cache) => cache.match(request.url));
  }

  /**
   * @param {Response} response
   * @returns {Boolean}
   */
  function cacheable(response) {
    return response.type === "basic" && response.ok && !response.headers.has("Content-Disposition")
  }

  event.respondWith(fetch(request).then(saveToCache).catch(serveFromCache));
});
