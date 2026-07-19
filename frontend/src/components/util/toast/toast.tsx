export interface IToast {
    message: string
    type: "error";
}

interface IToastProps {
    toast: IToast;
}

export default function Toast({
    toast,
}: IToastProps) {
    const bg = toast.type == 'error' ? 'bg-red-800' : '';
    return (
        <div className={`relative p-3 max-h-28 rounded-md ${bg} translate-x-full align-middle animate-fade-in-out-right opacity-0 z-10000`}>
            {toast.message}
        </div>
    )
}
