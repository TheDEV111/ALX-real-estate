import { propertyCategories } from "../data/propertyCategories";
import { useState } from "react";


function PropertyTypeFilter() {
  const [selected, setSelected] = useState(1);

  return (
    <div className="sticky top-18 h-[84px] flex overflow-x-auto gap-[43px] p-2 backdrop-blur-[54]   scrollbar-hide border-[#FDFDFD] border-b">
      {propertyCategories.map((category) => (
        <button
          key={category.id}
          onClick={() => setSelected(category.id)}
          className={`flex flex-col items-center gap-2 pb-2 whitespace-nowrap shrink-0 hover:scale-110 hover:brightness-110 transition-all duration-200 ${
            selected === category.id
              ? "border-b-2 border-gray-800"
              : "border-b-2 border-transparent"
          }`}
        >
          <div></div>
          <img src={category.icon} alt={category.name} className="w-6 h-6 " />
          <span className="text-[12px] text-[#616161] font-[500] ">
            {category.name}
          </span>
        </button>
      ))}
    </div>
  );
}

export default PropertyTypeFilter