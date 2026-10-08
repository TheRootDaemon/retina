import scrape from "../internal/retina/scripts/duckduckgo.js";

/**
 * Renders HTML into the document used by the scraper.
 *
 * @param {string} html - Search result page markup.
 */
const render = (html) => {
    document.body.innerHTML = html;
};

/**
 * Creates DuckDuckGo search result markup for a test case.
 *
 * @param {object} result - Search result fields.
 * @param {string|null} result.href - Result link, or null to omit it.
 * @param {string} result.title - Result title.
 * @param {string|null} result.snippet - Result snippet, or null to omit it.
 * @returns {string} HTML representing a DuckDuckGo result.
 */
const resultMarkup = ({ href, title, snippet }) => `
    <div class="result">
        ${href === null ? "" : `<a class="result__a" href="${href}">${title}</a>`}
        ${snippet === null ? "" : `<div class="result__snippet">${snippet}</div>`}
    </div>
`;

/**
 * Reset the document after each test to keep test cases isolated.
 */
afterEach(() => {
    document.body.innerHTML = "";
});

describe("duckduckgo scrape", () => {
    it("returns an empty array when there are no results", () => {
        render("<div>no results here</div>");

        expect(scrape()).toEqual([]);
    });

    it("returns an empty array when the result list is empty", () => {
        render("");

        expect(scrape()).toEqual([]);
    });

    it("extracts title, decoded url and snippet from a result", () => {
        render(
            resultMarkup({
                href: "https://duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.com%2F&rut=abc",
                title: "Example",
                snippet: "A snippet about example.",
            })
        );

        expect(scrape()).toEqual([
            {
                title: "Example",
                url: "https://example.com/",
                snippet: "A snippet about example.",
            },
        ]);
    });

    it("resolves protocol-relative result links", () => {
        render(
            resultMarkup({
                href: "//duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.com%2Fpath",
                title: "Example",
                snippet: "Snippet",
            })
        );

        expect(scrape()).toEqual([
            {
                title: "Example",
                url: "https://example.com/path",
                snippet: "Snippet",
            },
        ]);
    });

    it("trims surrounding whitespace from the title and snippet", () => {
        render(
            resultMarkup({
                href: "https://duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.com",
                title: "   Spaced Title   ",
                snippet: "\n   Spaced snippet   \n",
            })
        );

        expect(scrape()).toEqual([
            {
                title: "Spaced Title",
                url: "https://example.com",
                snippet: "Spaced snippet",
            },
        ]);
    });

    it("falls back to an empty url when there is no uddg parameter", () => {
        render(
            resultMarkup({
                href: "https://duckduckgo.com/about",
                title: "No redirect",
                snippet: "Snippet",
            })
        );

        expect(scrape()).toEqual([
            { title: "No redirect", url: "", snippet: "Snippet" },
        ]);
    });

    it("falls back to an empty snippet when there is none", () => {
        render(
            resultMarkup({
                href: "https://duckduckgo.com/l/?uddg=https%3A%2F%2Fexample.com",
                title: "No snippet",
                snippet: null,
            })
        );

        expect(scrape()).toEqual([
            { title: "No snippet", url: "https://example.com", snippet: "" },
        ]);
    });

    it("preserves the order of multiple results", () => {
        render(
            resultMarkup({
                href: "https://duckduckgo.com/l/?uddg=https%3A%2F%2Fa.example",
                title: "A",
                snippet: "first",
            }) +
                resultMarkup({
                    href: "https://duckduckgo.com/l/?uddg=https%3A%2F%2Fb.example",
                    title: "B",
                    snippet: "second",
                })
        );

        expect(scrape()).toEqual([
            { title: "A", url: "https://a.example", snippet: "first" },
            { title: "B", url: "https://b.example", snippet: "second" },
        ]);
    });

    it("falls back to an empty url when there is no result link", () => {
        render(
            resultMarkup({
                href: null,
                title: "No link",
                snippet: "Snippet",
            })
        );

        expect(scrape()).toEqual([
            {
                title: "",
                url: "",
                snippet: "Snippet",
            },
        ]);
    });
});
