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
  #manualResize() {
    this.#editContentAreaNode.style.height = "auto";
    this.#editContentAreaNode.style.height =
      this.#editContentAreaNode.scrollHeight + "px";
  }

  constructor(node) {
    this.#node = node;
    this.#editContentAreaNode = this.#node.querySelector(".edit-content-area");
    this.#previewContentAreaNode = this.#node.querySelector(
      ".preview-content-area .content"
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
   * This function converts #contentData to safe HTML,
   * parses all the allowed markdowns inside that HTML,
   * and then return the result.
   * @returns {string} safe HTML with some of the things inside it parsed
   */
  #parseHTMLSafely() {
    const lineStartMatches = [
      {
        matchStart: "### ",
        replaceStart: "<h3>",
        replaceEnd: "</h3>",
      },
      {
        matchStart: "## ",
        replaceStart: "<h2>",
        replaceEnd: "</h2>",
      },
      {
        matchStart: "# ",
        replaceStart: "<h1>",
        replaceEnd: "</h1>",
      },
      {
        matchStart: "#img:",
        replaceStart: '<img src="',
        replaceEnd: '">',
      },
      {
        matchStart: "#img-sm:",
        replaceStart: '<img class="small" src="',
        replaceEnd: '">',
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

    const tempElement = document.createElement("div");
    tempElement.innerText = this.#contentData;
    const safeHTML = tempElement.innerHTML.replaceAll("<br>", "\n");

    let result = "";

    const lines = safeHTML.split("\n");

    lines.forEach((line) => {
      let finalLine = "";
      let foundMatch = false;
      // Apply line starter matches, and <p>'s
      lineStartMatches.forEach((lineStartMatch) => {
        const lenMStart = lineStartMatch.matchStart.length;
        if (line.length >= lenMStart) {
          const startOfLine = line.substring(0, lenMStart);
          if (startOfLine === lineStartMatch.matchStart) {
            finalLine =
              lineStartMatch.replaceStart +
              line.substring(lenMStart) +
              lineStartMatch.replaceEnd;
            foundMatch = true;
          }
        }
      });
      if (!foundMatch) {
        finalLine = "<p>" + line + "</p>";
      }

      // Apply pair matches
      pairMatches.forEach((pairMatch) => {
        const countStart = finalLine.split(pairMatch.matchStart).length - 1;
        const equalMatchers = pairMatch.matchStart === pairMatch.matchEnd;
        if (equalMatchers) {
          const validMatches = Math.floor(countStart / 2);
          for (let i = validMatches; i > 0; i--) {
            finalLine = finalLine.replace(
              pairMatch.matchStart,
              pairMatch.replaceStart
            );
            finalLine = finalLine.replace(
              pairMatch.matchEnd,
              pairMatch.replaceEnd
            );
          }
        }
        if (!equalMatchers) {
          const countEnd = finalLine.split(pairMatch.matchEnd).length - 1;
          const validMatches = Math.min(countStart, countEnd);
          for (let i = validMatches; i > 0; i--) {
            finalLine = finalLine.replace(
              pairMatch.matchStart,
              pairMatch.replaceStart
            );
            finalLine = finalLine.replace(
              pairMatch.matchEnd,
              pairMatch.replaceEnd
            );
          }
        }
      });

      // add parsed line to the result string
      result += finalLine;
    });

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
      this.#previewContentAreaNode.innerHTML = this.#parseHTMLSafely();
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
    this.#manualResize();
  }
}
