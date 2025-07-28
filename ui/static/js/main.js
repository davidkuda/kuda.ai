console.log("Hello from the go server");

// ------------------------------------------------------------
// DOM references
// ------------------------------------------------------------
const themeBtn = document.getElementById("themeToggle");

// ------------------------------------------------------------
// event listeners
// ------------------------------------------------------------
themeBtn.addEventListener("click", async () => {
	const cur = document.documentElement.getAttribute("data-theme") || "dark";
	const next = cur === "dark" ? "light" : "dark";
	applyTheme(next);
	await localStorage.setItem("theme", next);
});

// HTMX Swaps:
// document.body.addEventListener('htmx:afterSwap', () => {
//   window.scrollTo({ top: 0});
// });
// document.body.addEventListener("htmx:afterSwap", () => {
// 	const main = document.querySelector("main");
// 	if (main) {
// 		main.scrollIntoView();
// 	}
// });
document.body.addEventListener("htmx:afterSwap", (e) => {
	const main = document.querySelector("main");
	// Only scroll when the swap target is <main>
	if (e.detail.target === main) {
		main.scrollIntoView();
	}
});


// ------------------------------------------------------------
// functions
// ------------------------------------------------------------
function applyTheme(mode) {
	document.documentElement.setAttribute("data-theme", mode);
	themeBtn.textContent = mode === "dark" ? "light mode" : "dark mode";
}

// ------------------------------------------------------------
// init
// ------------------------------------------------------------
(async () => {
	const theme = await localStorage.getItem("theme");
	applyTheme(theme);
})();
