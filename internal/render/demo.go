package render

import "time"

func Demo() Newsletter {
	return Newsletter{
		Brand:       "Tux Letter",
		Title:       "Supply-Chain Scares Meet a Quietly Stronger Linux Desktop",
		Subtitle:    "the AUR takes a hit while the kernel and Wayland keep maturing",
		Summary:     "A tense week for package trust on Arch collides with steady, meaningful progress across the kernel and the Wayland desktop.",
		GeneratedAt: time.Date(2026, 6, 25, 20, 0, 0, 0, time.UTC),
		Sections: []Section{
			{
				Heading: "package trust",
				Paragraphs: []string{
					"A short-lived set of malicious AUR packages renewed an old reminder: community repositories remain a real attack surface for advanced Linux users, and reviewing PKGBUILDs before building still matters [1].",
				},
			},
			{
				Heading: "kernel and desktop",
				Paragraphs: []string{
					"Linux 6.10 landed with broader hardware enablement and better power management, smoothing first-boot experiences on recent laptops [2].",
					"On the desktop, a long-running Wayland compositor latency series finally merged upstream, paying off years of incremental work and reducing perceived input lag [3].",
				},
			},
		},
		References: []Reference{
			{ID: 1, Title: "Active AUR malicious packages incident", Source: "Arch Linux", URL: "https://archlinux.org/news/active-aur-malicious-packages-incident"},
			{ID: 2, Title: "Linux Kernel 6.10 Released", Source: "9to5Linux", URL: "https://9to5linux.com/linux-kernel-6-10-released"},
			{ID: 3, Title: "Wayland Compositor Performance Series Merged", Source: "LWN", URL: "https://lwn.net/Articles/wayland-perf"},
		},
	}
}
