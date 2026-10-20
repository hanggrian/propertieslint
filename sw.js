const cacheName = self.location.pathname
const pages = [

  "/propertieslint/rules/comment/comment-spacing/",
    "/propertieslint/rules/comment/comment-style/",
    "/propertieslint/rules/comment/",
    "/propertieslint/rules/format/no-dangling-continuation/",
    "/propertieslint/rules/format/separator-style/",
    "/propertieslint/rules/format/valid-escape/",
    "/propertieslint/rules/format/",
    "/propertieslint/rules/pair/key-name/",
    "/propertieslint/rules/pair/no-duplicate-key/",
    "/propertieslint/rules/pair/no-empty-key/",
    "/propertieslint/rules/pair/no-empty-value/",
    "/propertieslint/rules/pair/",
    "/propertieslint/rules/whitespace/entry-trimming/",
    "/propertieslint/rules/whitespace/no-duplicate-blank-line/",
    "/propertieslint/rules/whitespace/no-leading-blank-line/",
    "/propertieslint/rules/whitespace/trailing-newline/",
    "/propertieslint/rules/whitespace/",
    "/propertieslint/",
    "/propertieslint/book/",
    "/propertieslint/categories/",
    "/propertieslint/rules/",
    "/propertieslint/tags/",
    "/propertieslint/book.min.98c83e4c8e1c8661368bcb29bbaf4dcc9d31b8164f76aae322d82d86a9750c65.css",
  "/propertieslint/en.search-data.min.4cb2dc4a6fad8deaa5e6c6149ffb0effee2f7a0f0cb9540c67360e664965debb.json",
  "/propertieslint/en.search.min.e4e21abe3f4c7d3f3cf1ce2525909e96a673d90cc816a761752bfbdfd24d624e.js",
  
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
