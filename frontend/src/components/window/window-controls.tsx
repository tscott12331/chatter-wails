import { Window } from "@wailsio/runtime"

export default function WindowControls() {
    const handleMaxMinimise = async () => {
        const isMaximized = await Window.IsMaximised();
        if(isMaximized) {
            Window.UnMaximise();
        } else {
            Window.Maximise();
        }
    }

    return (
        <div className="titlebar relative h-2 bg-bg-1 border-b border-b-outline-1 hover:*:h-8 hover:*:opacity-100 z-5000">
            <div className="absolute opacity-0 bg-inherit inset-0 w-full flex justify-end items-center text-center text-lg transition-[height_opacity]">
                <button onClick={Window.Minimise} className="w-8 h-8 rounded-none! border-none!">-</button>
                <button onClick={handleMaxMinimise} className="w-8 h-8 rounded-none! border-none!">□</button>
                <button onClick={Window.Close} className="w-8 h-8 rounded-none! border-none! hover:bg-red-500!">x</button>
            </div>
        </div>
    )
}
