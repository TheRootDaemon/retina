/**
 * Extracts search results from a test results page.
 *
 * This is intended to be used for testing retina.
 *
 * Each result is read from a `.result` element
 * and returns its title, destination URL, and snippet.
 *
 * Results without a link or snippet are returned with an empty string
 * for the corresponding field.
 *
 * @returns {Array<{
 *     title: string,
 *     url: string,
 *     snippet: string
 * }>} The extracted search results in page order.
 */
export default function scrape() {
    return Array.from(document.querySelectorAll(".result")).map((result) => {
        const link = result.querySelector(".result__a");

        return {
            title: link?.textContent?.trim() ?? "",
            url: link?.href ?? "",
            snippet:
                result.querySelector(".result__snippet")?.textContent?.trim() ??
                "",
        };
    });
}
