const form = document.getElementById("submit")

form.addEventListener('htmx:configRequest', e => {
	const year = e.detail.parameters.year
	const month = +e.detail.parameters.month
	const day = +e.detail.parameters.day

	const unix = moment(`${month}-${day}-${year}`, "jM-jD-jYYYY").unix()
	e.detail.parameters['datetimeEpoch'] = unix
})

