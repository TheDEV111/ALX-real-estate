import { useState } from "react";
import { MdOutlineFilterAlt } from "react-icons/md";
import { Button } from "./Button";
import { CustomSelect } from "./CustomSelect";
import { filterOptions, sortOptions } from "../data/propertyListingFilters";

const PropertyListingFilter = () => {
  const [selectedFilter, setSelectedFilter] = useState<number>(1);
  const [selectedSort, setSelectedSort] = useState<string>("Newest");
  const [isFilterOpen, setIsFilterOpen] = useState(false);


  return (
    <div className="flex place-content-between px-8 text-[12px] mb-[37px] font-semibold font-[Quicksand]">
      {/* Left side - Filter options */}
      <div className="flex gap-3 overflow-x-auto scrollbar-hide">
        {filterOptions.map((option) => (
          <Button
            key={option.id}
            onClick={() => setSelectedFilter(option.id)}
            className={`px-4 h-8 rounded-full transition-colors whitespace-nowrap  shrink-0 ${
              selectedFilter === option.id
                ? " text-[#34967C] bg-[#F0FFFB] border border-[#34967C]"
                : "border border-[#ECECEC] text-gray-800 hover:bg-gray-100"
            }`}
          >
            {option.name}
          </Button>
        ))}
      </div>

      {/* Right side - Filter and Sort buttons */}
      <div className="flex gap-2 ml-auto">
        <Button
          onClick={() => setIsFilterOpen(!isFilterOpen)}
          className={`flex items-center gap-1 px-4 h-8 rounded-full transition-colors ${
            isFilterOpen
              ? "text-[#34967C] border border-[#34967C]"
              : "border border-[#ECECEC] text-gray-800 hover:bg-gray-100"
          }`}
        >
          <MdOutlineFilterAlt className="w-5 h-5" />
          <span>Filter</span>
        </Button>

        <CustomSelect
          options={sortOptions}
          value={selectedSort}
          onChange={setSelectedSort}
        />
      </div>
    </div>
  );
};

export default PropertyListingFilter;
