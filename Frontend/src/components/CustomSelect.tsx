import { useState } from "react";

interface CustomSelectProps {
  options: string[];
  value: string;
  onChange: (value: string) => void;
  label?: string;
}

export const CustomSelect = ({
  options,
  value,
  onChange,
  label = "Sort by:",
}: CustomSelectProps) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="border border-[#ECECEC] px-4 h-8 rounded-full bg-white hover:bg-gray-50 flex items-center gap-1 transition-colors"
      >
        <span className="text-[#616161]">{label}</span>
        <span className="ml-1 text-black">{value}</span>
      </button>

      {isOpen && (
        <div className="absolute top-10 right-0 w-48 bg-white border border-[#ECECEC] rounded-lg shadow-lg z-10">
          {options.map((option) => (
            <button
              key={option}
              onClick={() => {
                onChange(option);
                setIsOpen(false);
              }}
              className={`w-full text-left px-4 py-2 text-[12px] transition-colors ${
                value === option
                  ? "bg-[#34967C] text-white"
                  : "text-[#616161] hover:bg-gray-50"
              }`}
            >
              {option}
            </button>
          ))}
        </div>
      )}
    </div>
  );
};
