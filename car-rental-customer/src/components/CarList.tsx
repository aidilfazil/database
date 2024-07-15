import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Box,
  Heading,
  SimpleGrid,
  Text,
  Button,
  VStack,
  useToast,
} from "@chakra-ui/react";
import axios from "axios";

interface Car {
  _id: string;
  make: string;
  model: string;
  year: number;
  type: string;
  available: boolean;
}

interface RentedCar {
  rental_id: string;
  car_id: string;
  make: string;
  model: string;
  year: number;
  type: string;
  rental_start_date: string;
  rental_end_date: string;
}

interface CarListProps {
  customerId: string;
}

const CarList: React.FC<CarListProps> = ({ customerId }) => {
  const queryClient = useQueryClient();
  const toast = useToast();
  const [currentCarId, setCurrentCarId] = useState<string | null>(null);

  const { data: cars, isLoading: carsLoading, error: carsError } = useQuery<Car[]>({
    queryKey: ['cars'],
    queryFn: () => axios.get("/api/cars").then((res) => res.data),
  });

  const { data: rentedCars, isLoading: rentedCarsLoading, error: rentedCarsError } = useQuery<RentedCar[]>({
    queryKey: ['rentedCars', customerId],
    queryFn: () => axios.get(`/api/rentals/customer/${customerId}`).then((res) => res.data),
  });

  const rentCarMutation = useMutation({
    mutationFn: (carId: string) => 
      axios.post("/api/rentals/create", { 
        car_id: carId,
        customer_id: customerId
      }),
    onMutate: (carId: string) => {
      setCurrentCarId(carId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ['cars']});
      queryClient.invalidateQueries({queryKey: ['rentedCars', customerId]});
      toast({
        title: "Car rented successfully",
        status: "success",
        duration: 3000,
        isClosable: true,
      });
      setCurrentCarId(null);
    },
    onError: () => {
      toast({
        title: "Failed to rent car",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
      setCurrentCarId(null);
    },
  });

  const returnCarMutation = useMutation({
    mutationFn: (rentalId: string) => 
      axios.post(`/api/rentals/${rentalId}/return`),
    onMutate: (rentalId: string) => {
      setCurrentCarId(rentalId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ['cars']});
      queryClient.invalidateQueries({queryKey: ['rentedCars', customerId]});
      toast({
        title: "Car returned successfully",
        status: "success",
        duration: 3000,
        isClosable: true,
      });
      setCurrentCarId(null);
    },
    onError: () => {
      toast({
        title: "Failed to return car",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
      setCurrentCarId(null);
    },
  });

  const handleRentNow = (carId: string) => {
    rentCarMutation.mutate(carId);
  };

  const handleReturnCar = (rentalId: string) => {
    returnCarMutation.mutate(rentalId);
  };

  if (carsLoading || rentedCarsLoading) return <Text>Loading...</Text>;
  if (carsError || rentedCarsError) return <Text>An error occurred</Text>;

  const availableCars = cars?.filter(car => car.available) || [];

  return (
    <Box width="100%">
      <Heading size="md" mb={4}>
        Available Cars
      </Heading>
      <SimpleGrid columns={[1, 2, 3]} spacing={4}>
        {availableCars.map((car) => (
          <Box key={car._id} borderWidth={1} borderRadius="lg" p={4}>
            <VStack align="start">
              <Text fontWeight="bold">
                {car.make} {car.model}
              </Text>
              <Text>Year: {car.year}</Text>
              <Text>Type: {car.type}</Text>
              <Button 
                colorScheme="green" 
                size="sm"
                onClick={() => handleRentNow(car._id)}
                isLoading={currentCarId === car._id && rentCarMutation.status === 'pending'}
              >
                Rent Now
              </Button>
            </VStack>
          </Box>
        ))}
      </SimpleGrid>

      <Heading size="md" mt={8} mb={4}>
        Your Rented Cars
      </Heading>
      <SimpleGrid columns={[1, 2, 3]} spacing={4}>
        {rentedCars?.map((rentedCar) => (
          <Box key={rentedCar.rental_id} borderWidth={1} borderRadius="lg" p={4}>
            <VStack align="start">
              <Text fontWeight="bold">
                {rentedCar.make} {rentedCar.model}
              </Text>
              <Text>Year: {rentedCar.year}</Text>
              <Text>Type: {rentedCar.type}</Text>
              <Text>Rented on: {new Date(rentedCar.rental_start_date).toLocaleDateString()}</Text>
              <Button 
                colorScheme="red" 
                size="sm"
                onClick={() => handleReturnCar(rentedCar.rental_id)}
                isLoading={currentCarId === rentedCar.rental_id && returnCarMutation.status === 'pending'}
              >
                Return Car
              </Button>
            </VStack>
          </Box>
        ))}
      </SimpleGrid>
    </Box>
  );
};

export default CarList;
