// Only tested with HTMX

export const $wrap = (node) => ({
  node: node,
  toggle: function (className, bool) {
    this.node.classList.toggle(className, bool);
    return this;
  },
  click: function (func) {
    this.node.addEventListener("click", func);
    return this;
  },
  getParent: function () {
    return $wrap(this.node.parentNode);
  },
  getHeight: function () {
    return this.node.scrollHeight;
  },
  getWidth: function () {
    return this.node.offsetWidth;
  },
});

export const $ = (query) => {
  const node = document.querySelector(query);
  return node === null ? null : $wrap(node);
};

export const $all = (query) => {
  return Array.from(document.querySelectorAll(query)).map((node) =>
    $wrap(node)
  );
};

export const $load = (func) => {
  window.addEventListener("load", func);

  // MDN about popstate:
  //
  // Note: When writing functions that process popstate event it is important
  // to take into account that properties like window.location will already reflect
  // the state change (if it affected the current URL), but document might still not.
  // If the goal is to catch the moment when the new document state is already fully in place,
  // a zero-delay setTimeout() method call should be used to effectively put its inner callback
  // function that does the processing at the end of the browser event loop.
  window.addEventListener("popstate", () => setTimeout(func, 0));
};

export const $htmxLoad = (func, executeNow) => {
  if (executeNow === true) func();
  window.addEventListener("htmx:load", func);
};
