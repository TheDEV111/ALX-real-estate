import { useState } from "react";
import { propertyCategories } from "../data/propertyCategories";


function PropertyTypeFilter()  {
  const [selected, setSelected] = useState<number>(1);

  return (
    <div className="sticky top-20 h-21 flex overflow-x-auto gap-10.75 py-4 px-8 backdrop-blur-[54px]   scrollbar-hide border-[#FDFDFD] border-b">
      {propertyCategories.map((category) => (
        <button
          key={category.id}
          onClick={() => setSelected(category.id)}
          className={`flex flex-col items-center gap-2 pb-14 whitespace-nowrap shrink-0 hover:scale-110 hover:brightness-110 hover:cursor-pointer transition-all duration-200 ${
            selected === category.id
              ? "border-b-2  border-gray-800"
              : "border-b-2 border-transparent"
          }`}
        >
          <img src={category.icon} alt={category.name} className="w-6 h-6 " />
          <span className="text-[12px]  text-[#616161] font-medium font-[Quicksand]">
            {category.name}
          </span>
        </button>
      ))}
    </div>
  );
}

export default PropertyTypeFilter;
