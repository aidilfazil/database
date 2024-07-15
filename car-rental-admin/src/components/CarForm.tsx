import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
  Select,
} from "@chakra-ui/react";
import axios from "axios";

const CarForm = () => {
  const [make, setMake] = useState("");
  const [model, setModel] = useState("");
  const [year, setYear] = useState("");
  const [type, setType] = useState("");

  const queryClient = useQueryClient();

  const createCarMutation = useMutation({
    mutationFn: (newCar: any) =>
      //axios.post("http://localhost:5000/api/cars", newCar),
      axios.post("/api/cars", newCar,{ withCredentials: true }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cars'] });
      setMake("");
      setModel("");
      setYear("");
      setType("");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createCarMutation.mutate({
      make,
      model,
      year: parseInt(year),
      type,
      available: true,
    });
  };

  return (
    <Box as="form" onSubmit={handleSubmit} width="100%">
      <VStack spacing={4}>
        <FormControl isRequired>
          <FormLabel>Make</FormLabel>
          <Input
            value={make}
            onChange={(e) => setMake(e.target.value)}
            placeholder="Enter car make"
          />
        </FormControl>
        <FormControl isRequired>
          <FormLabel>Model</FormLabel>
          <Input
            value={model}
            onChange={(e) => setModel(e.target.value)}
            placeholder="Enter car model"
          />
        </FormControl>
        <FormControl isRequired>
          <FormLabel>Year</FormLabel>
          <Input
            value={year}
            onChange={(e) => setYear(e.target.value)}
            placeholder="Enter car year"
            type="number"
          />
        </FormControl>
        <FormControl isRequired>
          <FormLabel>Type</FormLabel>
          <Select
            value={type}
            onChange={(e) => setType(e.target.value)}
            placeholder="Select car type"
          >
            <option value="Sedan">Sedan</option>
            <option value="SUV">SUV</option>
            <option value="Hatchback">Hatchback</option>
            <option value="Truck">Truck</option>
          </Select>
        </FormControl>
        <Button type="submit" colorScheme="blue">
          Add Car
        </Button>
      </VStack>
    </Box>
  );
};

export default CarForm;