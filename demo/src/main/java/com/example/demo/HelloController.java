package com.example.demo;

import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestMapping;

@RestController
public class HelloController {
    public static final String podName = System.getenv("POD_NAME");

    @RequestMapping("/")
    public String index() {
        return podName;
    }
}
