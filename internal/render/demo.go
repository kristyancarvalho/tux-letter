package render

import "time"

func Demo() Newsletter {
	return Newsletter{
		Title:       "Tux Letter",
		Summary:     "Kernel 6.10 lands with broad hardware support, Arch refreshes its installer, and a Wayland performance series finally merges. A calm but meaningful cycle for the open-source desktop.",
		GeneratedAt: time.Date(2026, 6, 25, 20, 0, 0, 0, time.UTC),
		Sources:     []string{"9to5linux.com", "archlinux.org", "lwn.net"},
		Items: []Item{
			{
				Title:        "Linux Kernel 6.10 Released",
				Source:       "9to5Linux",
				URL:          "https://9to5linux.com/linux-kernel-6-10-released",
				Summary:      "The 6.10 kernel ships with expanded hardware enablement, new filesystem features and improved power management for recent laptops.",
				WhyItMatters: "Better out-of-the-box hardware support means fewer post-install fixes on new machines.",
				Tags:         []string{"linux", "kernel", "release"},
			},
			{
				Title:        "Arch Linux Refreshes the Guided Installer",
				Source:       "Archlinux",
				URL:          "https://archlinux.org/news/installer-refresh",
				Summary:      "archinstall gains a cleaner flow and better disk handling, lowering the barrier for new Arch users.",
				WhyItMatters: "Makes a famously manual distro friendlier without compromising its philosophy.",
				Tags:         []string{"arch", "install"},
			},
			{
				Title:   "Wayland Compositor Performance Series Merged",
				Source:  "LWN",
				URL:     "https://lwn.net/Articles/wayland-perf",
				Summary: "A long-running patch series reducing compositor latency has finally been merged upstream.",
				Tags:    []string{"wayland", "graphics"},
			},
		},
	}
}
