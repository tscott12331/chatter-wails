import { ComponentPropsWithoutRef } from "react"

function SearchIcon(props: ComponentPropsWithoutRef<'svg'>) {
    return (
        <svg
            xmlns="http://www.w3.org/2000/svg"
            shapeRendering="geometricPrecision"
            textRendering="geometricPrecision"
            viewBox="0 0 300 300"
            {...props}
        >
            <g transform="matrix(1.08335 0 0 1.08335 0 0)">
                <path d="M0 100C0 44.772 44.772 0 100 0s100 44.772 100 100-44.772 100-100 100S0 155.228 0 100Zm100 43.1c23.803 0 43.1-19.297 43.1-43.1S123.803 56.9 100 56.9 56.9 76.197 56.9 100s19.297 43.1 43.1 43.1Z" />
                <rect
                    width={56.9}
                    height={30}
                    rx={0}
                    ry={0}
                    transform="rotate(-45 258.681 -66.915) scale(1 5.35524)"
                />
                <circle
                    r={43.1}
                    fill="rgba(0,0,0,0.13)"
                    transform="translate(100.1 100)"
                />
            </g>
        </svg>
    )
}

export default SearchIcon
