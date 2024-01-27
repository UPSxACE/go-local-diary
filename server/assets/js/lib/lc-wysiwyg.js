// TODO - Make it so it's visually clear when it's disabled

export class Editor {
  #contentData = "";
  #node = null;
  #editContentAreaNode = null;
  #previewContentAreaNode = null;
  #editButtonNode = null;
  #previewButtonNode = null;
  #wordCounterNode = null;
  #updatePreviewTimeout = null;

  constructor(node) {
    this.#node = node;
    this.#editContentAreaNode = this.#node.querySelector(".edit-content-area");
    this.#previewContentAreaNode = this.#node.querySelector(
      ".preview-content-area"
    );
    this.#editButtonNode = this.#node.querySelector('button[data-role="edit"]');
    this.#previewButtonNode = this.#node.querySelector(
      'button[data-role="preview"]'
    );
    this.#wordCounterNode = this.#node.querySelector(".word-counter .counter");

    const necessaryNodes = [
      this.#node,
      this.#editContentAreaNode,
      this.#previewContentAreaNode,
      this.#editButtonNode,
      this.#previewButtonNode,
      this.#wordCounterNode,
    ];
    const anyNullNode = necessaryNodes.some((node) => node === null);
    if (anyNullNode) {
      console.warn("Editor.constructor: One of the nodes was not found!");
      return;
    }
    this.#setupNodeEvents();
  }

  /**
   * This function needs to be called in the constructor, to add the necessary event listeners
   * to the HTML nodes of the Editor HTML.
   */
  #setupNodeEvents() {
    const onInput = function () {
      this.style.height = "auto";
      this.style.height = this.scrollHeight + "px";
    };

    this.#editContentAreaNode.setAttribute(
      "style",
      "height:" +
        this.#editContentAreaNode.scrollHeight +
        "px;overflow-y:hidden;"
    );

    this.#editContentAreaNode.addEventListener("input", onInput, false);

    this.#editContentAreaNode.addEventListener("input", (e) => {
      // NOTE - This implementation might change someday
      this.#contentData = e.target.value;
      this.#updatePreviewContent();
    });
    this.#editButtonNode.onclick = () => {
      this.#node.classList.toggle("preview", false);
      this.#updatePreviewContent();
    };
    this.#previewButtonNode.onclick = () => {
      this.#node.classList.toggle("preview", true);
      this.#updatePreviewContent();
    };
  }

  /**
   * This function is supposed to receive SAFE(sanitized) HTML as string argument,
   * ,parse all the allowed markdowns inside that HTML, and then return the result.
   * @param {string} safeHTML safe/sanitized HTML
   * @returns {string} HTML with some of the things inside it parsed
   */
  #parseSafeHTML(safeHTML) {
    let result = safeHTML;

    const lineFinishers = ["<br>", "</h1>", "</h2>", "</h3>"];

    const lineStartMatches = [
      {
        matchStart: "### ",
        matchEnd: "<br>",
        replaceStart: "<h3>",
        replaceEnd: "</h3>",
      },
      {
        matchStart: "## ",
        matchEnd: "<br>",
        replaceStart: "<h2>",
        replaceEnd: "</h2>",
      },
      {
        matchStart: "# ",
        matchEnd: "<br>",
        replaceStart: "<h1>",
        replaceEnd: "</h1>",
      },
    ];

    const pairMatches = [
      {
        matchStart: "**",
        matchEnd: "**",
        replaceStart: "<strong>",
        replaceEnd: "</strong>",
      },
      {
        matchStart: "*",
        matchEnd: "*",
        replaceStart: "<em>",
        replaceEnd: "</em>",
      },
      // {
      //   matchStart: "&lt;h1&gt;",
      //   matchEnd: "&lt;/h1&gt;",
      //   replaceStart: "<h1>",
      //   replaceEnd: "</h1>",
      // },
      // {
      //   matchStart: "&lt;h1&gt;",
      //   matchEnd: "&lt;/h1&gt;<br>",
      //   replaceStart: "<h1>",
      //   replaceEnd: "</h1>",
      // },
    ];

    while (lineStartMatches.length > 0 || pairMatches.length > 0) {
      if (lineStartMatches.length > 0) {
        const nextMatch = lineStartMatches.shift();
        let noneFound = false;
        let ignoreIndexes = -1;
        while (!noneFound) {
          const initialPosition = ignoreIndexes + 1;

          const indexStart = result.indexOf(
            nextMatch.matchStart,
            initialPosition
          );
          const indexEnd = result.indexOf(nextMatch.matchEnd, indexStart + 1);

          const followsBreak = lineFinishers.some((lineFinisher) => {
            const len = lineFinisher.length;
            if (indexStart < len) return false;

            return (
              result.substring(indexStart - len).indexOf(lineFinisher) === 0
            );
          });

          let lineStart = indexStart === 0 || followsBreak;

          if (lineStart) {
            if (indexStart !== -1 && indexEnd !== -1) {
              let temp1 = result.substring(0, indexStart);
              let temp2 = result.substring(indexStart);

              temp2 = temp2.replace(
                nextMatch.matchStart,
                nextMatch.replaceStart
              );
              temp2 = temp2.replace(nextMatch.matchEnd, nextMatch.replaceEnd);

              result = temp1 + temp2;
            }
            if (indexStart !== -1 && indexEnd === -1) {
              let temp1 = result.substring(0, indexStart);
              let temp2 = result.substring(indexStart);

              temp2 = temp2.replace(
                nextMatch.matchStart,
                nextMatch.replaceStart
              );
              temp2 += nextMatch.replaceEnd;

              result = temp1 + temp2;
            }
          }

          if (!lineStart) {
            ignoreIndexes = indexStart;
          }

          if (indexStart === -1) {
            noneFound = true;
          }
        }
      }

      if (pairMatches.length > 0) {
        const nextMatch = pairMatches.shift();
        let noneFound = false;

        while (!noneFound) {
          const indexStart = result.indexOf(nextMatch.matchStart);
          const indexEnd = result.indexOf(nextMatch.matchEnd, indexStart + 1);
          if (indexStart !== -1 && indexEnd !== -1) {
            let temp1 = result.substring(0, indexStart);
            let temp2 = result.substring(indexStart);

            temp2 = temp2.replace(nextMatch.matchStart, nextMatch.replaceStart);
            temp2 = temp2.replace(nextMatch.matchEnd, nextMatch.replaceEnd);

            result = temp1 + temp2;
          }
          if (indexStart === -1 || indexEnd === -1) {
            noneFound = true;
          }
        }
      }
    }

    return result;
  }

  /**
   * This function must be called to update the parsed content in the preview mode.
   * It must be called whenever the content changes.
   * It also is responsible to call the updateWordCount function.
   */
  #updatePreviewContent() {
    this.#updatePreviewTimeout = clearTimeout(this.#updatePreviewTimeout);
    this.#updatePreviewTimeout = setTimeout(() => {
      // NOTE - This implementation might change someday
      this.#previewContentAreaNode.innerText = this.#contentData;
      this.#previewContentAreaNode.innerHTML = this.#parseSafeHTML(
        this.#previewContentAreaNode.innerHTML
      );
      this.#updateWordCount();
    }, 500);
  }

  /**
   * This function must be called whenever there is a change in the content,
   * to update the word counter.
   */
  #updateWordCount() {
    // TODO - Replace this by a regular expression after the rest is correctly implemented
    // and make it so it doesn't count symbols as words, or double spaces, or markdown elements
    const cleanWordsArray = this.#previewContentAreaNode.innerText
      .replaceAll("\n", " ") // replace line breaks by spaces
      .trim() // remove trailing spaces on end and begin
      .split(" ") // divide string in array of words
      .filter((x) => x !== ""); // remove empty words caused by sequential spaces on the text

    this.#wordCounterNode.innerHTML = cleanWordsArray.length;
  }

  /** Allows the user to type on the editor. */
  enable() {
    this.#editContentAreaNode.removeAttribute("disabled");
  }

  /** Blocks the user the user from typing. */
  disable() {
    this.#editContentAreaNode.setAttribute("disabled");
  }

  getContent() {
    return this.#contentData;
  }

  setContent(newContent) {
    this.#contentData = newContent;
    this.#editContentAreaNode.value = newContent;
    this.#updatePreviewContent();
  }
}
