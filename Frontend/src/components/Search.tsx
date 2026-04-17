import { BiSearch } from "react-icons/bi";


function Search() {
  return (
    <div className="border-[1px] w-full md:w-auto py-2 rounded-full shadow-sm hover:shadow-md transition cursor-pointer">
      <div className="flex flex-row items-center justify-between gap-5">
        <div className="text-sm px-6 border-r-[1px]">
          Location
        </div>
        <div className="hidden sm:block text-sm px-6 border-r-[1px] flex-1 text-center">
          Check in
        </div>
        <div className="hidden sm:block text-sm  px-6 flex-1 border-r-[1px] text-center w-max">
          Check out
        </div>
        <div className="text-xs pl-6 pr-2 text-gray-600 flex flex-row items-center gap-6">
          <div className="hidden sm:block">Add guests</div>
          <div className="p-2 bg-orange-500 rounded-full text-white">
            <BiSearch size={18} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default Search;
