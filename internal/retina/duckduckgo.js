Array.from(document.querySelectorAll(".result")).map((result) => {
    const link = result.querySelector(".result__a");

    const url = new URL(link?.href ?? "").searchParams.get("uddg");

    return {
        title: link?.textContent?.trim() ?? "",
        url: url ?? "",
        snippet:
            result.querySelector(".result__snippet")?.textContent?.trim() ?? "",
    };
});
