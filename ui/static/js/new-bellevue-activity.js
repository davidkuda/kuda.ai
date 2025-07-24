document.querySelectorAll(".number-picker").forEach((picker) => {
	const input = picker.querySelector("input");
	picker
		.querySelector(".minus")
		.addEventListener("click", () => input.stepDown());
	picker.querySelector(".plus").addEventListener("click", () => input.stepUp());
});
