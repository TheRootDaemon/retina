/**
 * Extracts search results from a DuckDuckGo results page.
 *
 * Each result is read from a `.result` element
 * and returns its title, destination URL, and snippet.
 * DuckDuckGo wraps result URLs in a redirect link
 * containing the destination in the `uddg` query parameter.
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

        const url = link
            ? new URL(link?.href ?? "").searchParams.get("uddg")
            : null;

        return {
            title: link?.textContent?.trim() ?? "",
            url: url ?? "",
            snippet:
                result.querySelector(".result__snippet")?.textContent?.trim() ??
                "",
        };
    });
}
