import * as htmx from "htmx.org";
import * as feather from "feather-icons/dist/feather.min";
import Alpine from "alpinejs";
import * as noteEditor from "./editor";

// $state will save some persistent data so it's not lost on htmx-boost page changes (related to UI mostly)
window.$state = {};
window.htmx = htmx;
window.Alpine = Alpine;
window.noteEditor = noteEditor;

// feather.js
$load(() => {
  feather.replace();
  $htmxLoad(() => feather.replace());
});

Alpine.start();

// Add class .loaded to the <body> element after page is fully loaded, to easily add page entering transitions
$realLoad(() => {
  document.body.classList.add("loaded");
});

// Force reload when using the go back page button (history API) to avoid javascript related bugs
// window.addEventListener("popstate", () => location.reload());
