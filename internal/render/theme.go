package render

type Theme struct {
	Background string
	Surface    string
	Surface2   string
	Text       string
	Muted      string
	Accent     string
	Accent2    string
	Warning    string
	Border     string
}

func Cypherpunk() Theme {
	return Theme{
		Background: "#0b0f14",
		Surface:    "#111821",
		Surface2:   "#17212b",
		Text:       "#e6edf3",
		Muted:      "#9aa7b2",
		Accent:     "#00e5ff",
		Accent2:    "#8cff66",
		Warning:    "#ffcc66",
		Border:     "#263442",
	}
}

const (
	fontBody = "-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif"
	fontMono = "'SFMono-Regular',ui-monospace,Menlo,Consolas,'Liberation Mono',monospace"
)
