import { MdBed, MdBathtub, MdPeople } from "react-icons/md";

import { useState } from "react";
import type { PropertyListing } from "../data/propertyListings";

interface PropertyCardProps {
  property: PropertyListing;
}

export const PropertyCard: React.FC<PropertyCardProps> = ({ property }) => {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  return (
    <div className="bg-white rounded-xl overflow-hidden cursor-pointer mb-7 pb-4 font-[Quicksand]">
      {/* Image Container */}
      <div className="relative w-full max-w-[379px] md:max-w-[454px] h-[264px] md:h-[305px] lg:h-[255px] bg-gray-200 overflow-hidden rounded-xl">
        <img
          src={property.image}
          alt={property.name}
          className="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
        />

        {/* Discount Badge */}

        {property.discount && (
          <div className="flex gap-1 items-center absolute top-3 left-0 bg-[#34967C] text-white px-3 py-3 rounded-r-full rounded-tl-full text-[15px] ">
            <img src="/icons/Subtract.svg" alt="discount" className="w-4 h-4" />
            {property.discount}
          </div>
        )}

        {/* Image Carousel Dots */}
        <div className="absolute bottom-3 left-1/2 transform -translate-x-1/2 flex gap-1">
          {[0, 1, 2].map((index) => (
            <button
              key={index}
              onClick={() => setCurrentImageIndex(index)}
              className={`w-2 h-2 rounded-full transition ${
                currentImageIndex === index ? "bg-white" : "bg-white/50"
              }`}
            />
          ))}
        </div>
      </div>

      {/* Content */}
      <div className="pt-4">
        {/* Categories */}
        <div className="flex flex-nowrap mb-3 gap-2  ">
          {property.categories.map((category) => (
            <span
              key={category}
              className="text-[10px] md:text-[12px] px-2 py-1 bg-gray-100 text-gray-700 rounded-full font-medium whitespace-nowrap "
            >
              {category}
            </span>
          ))}
        </div>

        {/* Name and Rating */}
        <div className="flex justify-between items-start mb-1">
          <h2 className="font-semibold text-[16px] md:text-[18px] text-gray-900 font-[Quicksand] whitespace-nowrap">
            {property.name}
          </h2>
          <div className="flex items-center gap-1.5 ml-2">
            <span className="text-[#FAC02B] text-[18px]">★</span>
            <span className="text-[14px] font-medium text-gray-900">
              {property.rating}
            </span>
          </div>
        </div>

        {/* Location */}
        <p className="text-[12px] md:text-[14px]  mb-3 font-[Quicksand]">
          {property.location}
        </p>

        {/* Amenities and Price*/}
        <div className="flex items-center justify-between  mt-4 ">
          {/* Amenities */}
          <div className=" flex gap-2 border border-gray-200 rounded-3xl px-2.5 py-1.5 text-[8px]">
            <div className=" flex items-center gap-1 text-[12px]">
              <MdBed size={16} />
              <span>{property.beds}</span>
            </div>
            <div className="flex items-center gap-1 text-[12px]">
              <MdBathtub size={16} />
              <span>{property.bathrooms}</span>
            </div>
            <div className="flex items-center gap-1 text-[12px]">
              <MdPeople size={16} />
              <span>{property.capacity}</span>
            </div>
          </div>

          {/* Price */}
          <p className="text-[13px] md:text-[14px] text-gray-900">
            <span className="font-semibold text-[22px] md:text-[17px]">
              ${property.price}
            </span>
            <span className="text-gray-600">/n</span>
          </p>
        </div>
      </div>
    </div>
  );
};
